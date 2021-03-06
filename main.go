package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"

	"github.com/young-steveo/godax/gdax"
	"github.com/young-steveo/godax/market"

	"github.com/buger/jsonparser"
	"github.com/joho/godotenv"
)

func main() {
	var leftTicker string
	var rightTicker string

	flag.StringVar(&leftTicker, `left`, `LTC`, `Left side of pair, e.g. LTC`)
	flag.StringVar(&rightTicker, `right`, `BTC`, `Right side of pair, e.g. BTC`)
	flag.Parse()

	productID := market.GetProductID(leftTicker, rightTicker)

	log.SetOutput(os.Stdout)
	log.Println(`¡Hola! Let's make some trades`)
	log.Println(`===`)

	log.Println(`Bootstrapping environment`)

	err := godotenv.Load()
	if err != nil {
		log.Fatal(`Error loading .env file`, err.Error())
	}

	/**
	 * Signals
	 */
	sigs := make(chan os.Signal)
	exit := make(chan int)
	stopFill := make(chan int)
	signal.Notify(sigs, os.Interrupt)

	/**
	 * Message Channels
	 */
	done := make(chan []byte)
	received := make(chan []byte)

	/**
	 * Connection
	 */
	err = gdax.Connect()
	if err != nil {
		log.Fatal("Could not connect: ", err)
		return
	}
	defer gdax.Close() // close the websocket when the main function exits.

	/**
	 * Consume all websocket messages and route them to the proper channel.
	 */
	go func() {
		defer close(exit) // if this goroutine returns, close the exit channel
		for {
			_, message, err := gdax.ReadMessage()
			if err != nil {
				log.Println("Could not read from socket: ", err)
				gdax.Unsubscribe()
				return
			}

			typ, _ := jsonparser.GetUnsafeString(message, `type`)

			switch typ {
			case `received`:
				received <- message
			case `done`:
				done <- message
			}
		}
	}()

	/**
	 * Received messages from GDAX
	 */
	go func() {
		orders := make(chan *market.Order)
		defer close(orders)
		for {
			if message, more := <-received; more {
				// need to match client order id with server order id and update the order
				coid, err := jsonparser.GetUnsafeString(message, `client_oid`)
				if err != nil {
					log.Printf(`Could not read client_oid`)
					continue
				}
				clientID, err := uuid.Parse(coid)
				if err != nil {
					log.Printf(`Could not parse client_oid`)
					continue
				}

				oid, err := jsonparser.GetUnsafeString(message, `order_id`)
				if err != nil {
					log.Printf(`Could not read order_id`)
					continue
				}
				orderID, err := uuid.Parse(oid)
				if err != nil {
					log.Printf(`Could not parse order_id`)
					continue
				}

				gdax.SyncOrderID(clientID, orderID)
			} else {
				break
			}
		}
	}()

	/**
	 * Done messages from GDAX
	 */
	go func() {
		for {
			if msg, more := <-done; more {
				// need to clean out local order and place new orders
				log.Println(string(msg))
				reason, err := jsonparser.GetString(msg, `reason`)
				if err != nil {
					log.Println(`Error reading reason from done message.  Skipping.`)
					continue
				}
				if reason == `canceled` {
					log.Println(`Order canceled, skipping done handler.`)
					continue
				}
				oid, err := jsonparser.GetString(msg, `order_id`)
				if err != nil {
					log.Println(`Error reading order_id from done message.  Skipping.`)
					continue
				}
				orderID, err := uuid.Parse(oid)
				if err != nil {
					log.Println(`Error parsing order_id from done message.  Skipping.`)
					continue
				}
				gdax.PlaceSpreadFrom(orderID)
				gdax.RemoveOrder(orderID)
			} else {
				break
			}
		}
	}()

	/**
	 * Fill gaps
	 */
	go func() {
		defer close(stopFill)
		for {
			select {
			case <-stopFill:
				break
			default:
				gdax.FillGaps(productID)
				time.Sleep(5 * time.Second)
			}
		}
	}()

	/**
	 * Tracking
	 */
	var start float64
	var end float64

	/**
	 * Graceful shutdown.  Unsubscribes to the websocket
	 */
	go func() {
		defer close(exit) // close the exit channel when we are finished cleaning up.
		<-sigs            // wait for a signal.
		log.Printf("Shutting down")
		stopFill <- 0
		gdax.Unsubscribe()
		end, _, _ = gdax.GetBalance(productID)
		log.Printf(`Delta in BTC: %f`, end-start)
	}()

	/**
	 * Subscribe message to GDAX initiates the websocket messages.
	 */
	gdax.Subscribe(productID)

	accts, err := gdax.GetAccounts()
	if err != nil {
		log.Fatal(`error: ` + err.Error())
	}

	acct := accts.GetByCurrency(productID[0])

	

	log.Println(`Canceling all pending orders`)
	_, err = gdax.Request(`DELETE`, `/orders`, nil)
	if err != nil {
		log.Fatal(`Cancel Order request failed: ` + err.Error())
	}

	// wait for rebalance.
	var orderStatus int
	orderStatus, start, err = gdax.Rebalance(productID)

	if err != nil {
		log.Fatal(`Rebalance Failed: ` + err.Error())
	}
	switch orderStatus {
	case -1:
		log.Fatal(`error placing rebalance order`)
	case 0:
		log.Println(`No rebalance was needed.`)
		_, price, err := gdax.GetBalance(productID)
		if err != nil {
			log.Fatal(`GetBalance Failed: ` + err.Error())
		}
		gdax.PlaceSpread(productID, price)
	case 1:
		log.Println(`Rebalance order placed.  Once the done signal sends, we can place our spread.`)
	default:
		log.Fatal(`Unknown order satatus code.  You done f'ed up your code.`)
	}

	/**
	 * Block until an exit message is received
	 */
	exitCode := <-exit
	fmt.Println("Bye!")
	os.Exit(exitCode)
}
