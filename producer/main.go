package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

const queueName = "Service1Queue"

func main() {
	// new connection
	conn, err := amqp.Dial("amqp://localhost:5672")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// open channel
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	// declare the queue
	q, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		panic(err)
	}

	// Create a new Gin instance
	router := gin.Default()

	// Add route to send a message to the queue
	router.GET("/send", func(ctx *gin.Context) {
		msg := ctx.Query("msg")
		if msg == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Message is required"})
			return
		}

		// Create a message to publish
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		}

		// Publish a message to the queue
		err = channel.Publish(
			"",      // exchange
			q.Name,  // routing key
			false,   // mandatory
			false,   // immediate
			message, // actual message
		)
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish the message"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": msg, "status": "success"})
	})

	// Start the Gin server
	log.Fatal(router.Run(":8080"))
}
