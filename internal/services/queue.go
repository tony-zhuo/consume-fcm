package services

import (
	"context"
	"encoding/json"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/tony-zhuo/consume-fcm/internal/models"
	"github.com/tony-zhuo/consume-fcm/internal/repositories"
	"time"
)

type QueueService interface {
	Producer(data MessageQueue) error
	Consumer()
}

type queueService struct {
	queueConn  *amqp.Connection
	fcmClient  *messaging.Client
	fcmJobRepo repositories.FcmJobRepo
}

func NewQueueService(
	queueConn *amqp.Connection,
	fcmClient *messaging.Client,
	fcmJobRepo repositories.FcmJobRepo) QueueService {

	return &queueService{
		queueConn:  queueConn,
		fcmClient:  fcmClient,
		fcmJobRepo: fcmJobRepo,
	}
}

type MessageQueue struct {
	Identifier string `json:"identifier" validate:"required"`
	Type       string `json:"type" validate:"required"`
	DeviceId   string `json:"device_id" validate:"required"`
	Text       string `json:"text" validate:"required"`
}

// Producer add message queue to notification.fcm
func (s *queueService) Producer(data MessageQueue) error {
	channelRabbitMQ, err := s.queueConn.Channel()
	if err != nil {
		fmt.Println("queue conn channel error: ", err)
		return err
	}
	defer channelRabbitMQ.Close()

	// convert struct to []byte
	body, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	message := amqp.Publishing{
		MessageId:   uuid.New().String(),
		ContentType: "application/json",
		Body:        body,
	}

	if err := channelRabbitMQ.Publish(
		"",
		"notification.fcm",
		false,
		false,
		message,
	); err != nil {
		fmt.Println("queue publish error: ", err)
		return err
	}

	return nil
}

// Consumer get the message from notification.fcm
func (s *queueService) Consumer() {
	channelRabbitMQ, err := s.queueConn.Channel()
	if err != nil {
		return
	}
	defer channelRabbitMQ.Close()

	q, err := channelRabbitMQ.QueueDeclare(
		"notification.fcm",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return
	}

	messages, err := channelRabbitMQ.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return
	}

	for message := range messages {
		var messageQueue MessageQueue
		if err := json.Unmarshal(message.Body, &messageQueue); err != nil {
			fmt.Println(err)
		}

		// validate
		if !s.validator(messageQueue) {
			// if invalid fields, then skip
			fmt.Println("validator fail")
			continue
		}

		// push notification
		_, err := s.pushNotification(context.Background(), messageQueue.Text)
		if err != nil {
			// if push notification fail, then skip
			fmt.Println("pushNotification err: ", err)
			continue
		}

		// store to database
		job, err := s.storeNotification(messageQueue.Identifier)
		if err != nil {
			fmt.Println("storeNotification err: ", err)
		}

		// publish message
		if err := s.publish(*job); err != nil {
			fmt.Println("storeNotification err: ", err)
		}
	}

	return
}

func (s *queueService) validator(messageQueue MessageQueue) bool {
	validate := validator.New()
	err := validate.Struct(messageQueue)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (s *queueService) pushNotification(ctx context.Context, text string) (string, error) {
	return s.fcmClient.Send(ctx, &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Incoming message",
			Body:  text,
		},
		Topic: "test",
	})
}

func (s *queueService) storeNotification(identifier string) (*models.FcmJob, error) {
	job := models.FcmJob{
		Identifier: identifier,
		DeliverAt:  time.Now(),
	}
	if err := s.fcmJobRepo.Create(job); err != nil {
		return nil, err
	}
	return &job, nil
}

func (s *queueService) publish(data models.FcmJob) error {
	fmt.Println("publish data", data)
	channelRabbitMQ, err := s.queueConn.Channel()
	if err != nil {
		return err
	}
	defer channelRabbitMQ.Close()

	body, err := json.Marshal(data)
	if err != nil {
		fmt.Println("publish json marshal err: ", err)
		return err
	}

	if err := channelRabbitMQ.ExchangeDeclare(
		"notification.done",
		"topic",
		true,
		false,
		false,
		false,
		nil); err != nil {

		fmt.Println("publish json marshal err: ", err)
		return err
	}

	message := amqp.Publishing{
		MessageId:   uuid.New().String(),
		ContentType: "application/json",
		Body:        body,
	}

	if err := channelRabbitMQ.Publish(
		"",
		"notification.done",
		false,
		false,
		message,
	); err != nil {
		fmt.Println("publish error", err)
		return err
	}

	return nil
}
