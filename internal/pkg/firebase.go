package pkg

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"google.golang.org/api/option"
	"os"
)

func NewFirebaseConnection() *messaging.Client {
	opt := option.WithCredentialsFile("test-firebase-account-key.json")
	projectId := os.Getenv("FIREBASE_PROJECT_ID")

	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: projectId,
	}, opt)
	if err != nil {
		panic(fmt.Errorf("error initializing app: %v", err))
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		panic(fmt.Errorf("error initializing app: %v", err))
	}

	return client
}
