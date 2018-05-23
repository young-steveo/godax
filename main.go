package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/young-steveo/godax/market"

	"github.com/young-steveo/godax/message"

	"github.com/buger/jsonparser"
	"github.com/joho/godotenv"
	"github.com/young-steveo/godax/gdax"
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

	gdax.PlaceOrder(market.MakeOrder("sell", "1.0", "127.12"))

	exitCode := <-exit // wait for something to close.
	fmt.Println("Bye!")
	os.Exit(exitCode)
}
