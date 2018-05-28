package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/google/uuid"

	"github.com/young-steveo/godax/gdax"
	"github.com/young-steveo/godax/market"
	"github.com/young-steveo/godax/message"

	"github.com/buger/jsonparser"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println(string(message.Subscribe()))
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
	 * Manage the Book
	 */
	orderCMDs := make(chan market.OrderCommand)
	accountCMDs := make(chan market.AccountCommand)
	go func() {
		close := make(chan bool)
		keeper := market.MakeKeeper(orderCMDs, accountCMDs)
		keeper.Listen(close) // contains a blocking loop
	}()

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
			log.Println(string(message))

			typ, _ := jsonparser.GetUnsafeString(message, `type`)

			switch typ {
			case `done`:
				done <- message
			case `received`:
				received <- message
			}
		}
	}()

	/**
	 * Done messages from GDAX
	 */
	go func() {
		for {
			if _, more := <-done; more {

				// need to clean out local order and place new orders
			} else {
				break
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

				orderCMDs <- &market.UpdateOrder{ClientIDValue: clientID, ServerIDValue: orderID, Resp: orders}
			} else {
				break
			}
		}
	}()

	/**
	 * Graceful shutdown.  Unsubscribes to the websocket
	 */
	go func() {
		defer close(exit) // close the exit channel when we are finished cleaning up.
		<-sigs            // wait for a signal.
		log.Printf("Shutting down")
		gdax.Unsubscribe()
	}()

	/**
	 * Subscribe message to GDAX initiates the websocket messages.
	 */
	//gdax.Subscribe()

	// TEMPORARY manual code to make an order
	// o := market.MakeOrder(market.Side("sell"), market.Size("5.0"), market.Price("117.03"))

	// orders := make(chan *market.Order)
	// defer close(orders)
	// create := &market.CreateOrder{Order: o, Resp: orders}
	// orderCMDs <- create
	// create.Orders() <- o

	// log.Printf(`Placing an order to %s %s LTC for $%s`, o.Side, o.Size, o.Price)
	// gdax.Request(`POST`, `/orders`, o)

	accounts, err := gdax.GetAccounts()
	if err != nil {
		log.Fatal(`error: ` + err.Error())
	}
	for _, account := range accounts {
		accountChannel := make(chan *market.Account)
		add := &market.AddAccount{Ticker: account.Currency, Resp: accountChannel}
		accountCMDs <- add
		add.Accounts() <- account
		<-add.Accounts() // wait till it's done
	}

	log.Println(`Canceling all pending orders`)
	_, err = gdax.Request(`DELETE`, `/orders`, nil)
	if err != nil {
		log.Fatal(`Cancel Order request failed: ` + err.Error())
	}

	<-gdax.Rebalance(market.GetProductID(`LTC`, `BTC`), accountCMDs, orderCMDs) // wait for rebalance.

	/**
	 * Block until an exit message is received
	 */
	exitCode := <-exit
	fmt.Println("Bye!")
	os.Exit(exitCode)
}
