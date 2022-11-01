package pkg

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/streadway/amqp"
	"os"
)

func NewRabbitmqConnection() *amqp.Connection {
	user := os.Getenv("RABBITMQ_USER")
	pwd := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")

	conn := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pwd, host, port)

	if conn, err := amqp.Dial(conn); err != nil {
		panic("new RabbitMQ error")
	} else {
		return conn
	}
}
