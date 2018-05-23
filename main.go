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
	 * Channels
	 */
	sigs := make(chan os.Signal)
	exit := make(chan int)
	done := make(chan []byte)
	received := make(chan []byte)

	signal.Notify(sigs, os.Interrupt)

	/**
	 * Connection
	 */
	err = gdax.Connect()
	if err != nil {
		log.Fatal("Could not connect: ", err)
		return
	}
	defer gdax.Close() // close the websocket when the main function exits.

	orders := make(market.MyBook, 0)

	// Consume messages
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
			case `done`:
				done <- message
			case `received`:
				received <- message
			}
		}
	}()

	// Handle Done messages
	go func() {
		for {
			if message, more := <-done; more {
				log.Println(string(message))
				// need to clean out local order and place new orders
			} else {
				break
			}
		}
	}()

	// Handle Received messages
	go func() {
		for {
			if message, more := <-received; more {
				log.Println(string(message))
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

				order := orders.GetByClientID(clientID)
				if order == nil {
					log.Printf(`Could not find order by client_oid: %s`, coid)
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

				order.ServerID = orderID
			} else {
				break
			}
		}
	}()

	// Graceful shut down
	go func() {
		defer close(exit) // close the exit channel when we are finished cleaning up.
		<-sigs            // wait for a signal.
		log.Printf("Shutting down")
		gdax.Unsubscribe()
	}()

	// Send Subscribe message to GDAX
	gdax.Subscribe()

	o := market.MakeOrder("sell", "1.0", "127.13")

	orders = append(orders, o)

	log.Printf(`Placing an order to %s %s LTC for $%s`, o.Side, o.Size, o.Price)
	gdax.Request(`POST`, `/orders`, o)

	exitCode := <-exit // wait for something to close.
	fmt.Println("Bye!")
	os.Exit(exitCode)
}
