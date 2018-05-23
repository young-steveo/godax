package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/young-steveo/godax/gdax"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println(`¡Hola! Let's make some trades`)
	log.Println(`===`)

	log.Println(`Bootstrapping environment`)

	err := godotenv.Load()
	if err != nil {
		log.Fatal(`Error loading .env file`, err.Error())
	}

	/**
	 * Signal Processing
	 */
	sigs := make(chan os.Signal)
	done := make(chan int)

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
		defer close(done) // if this goroutine returns, close the done channel
		for {
			_, message, err := gdax.ReadMessage()
			if err != nil {
				log.Println("Could not read from socket: ", err)
				gdax.Unsubscribe()
				return
			}
			log.Println(string(message))
		}
	}()

	// Graceful shut down
	go func() {
		defer close(done) // close the done channel when we are finished cleaning up.
		<-sigs            // wait for a signal.
		log.Printf("Shutting down")
		gdax.Unsubscribe()
	}()

	// Send Subscribe message to GDAX
	gdax.Subscribe()

	exitCode := <-done // wait for something to close.
	fmt.Println("Bye!")
	os.Exit(exitCode)
}
