package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tony-zhuo/consume-fcm/internal/controllers"
	"github.com/tony-zhuo/consume-fcm/internal/pkg"
	"github.com/tony-zhuo/consume-fcm/internal/repositories"
	"github.com/tony-zhuo/consume-fcm/internal/services"
	"github.com/tony-zhuo/consume-fcm/internal/workers"
)

func main() {
	// new RabbitMQ connection
	queueConn := pkg.NewRabbitmqConnection()
	// new Firebase connection
	firebaseClient := pkg.NewFirebaseConnection()
	// new database connection
	database := pkg.NewDatabaseConnection()

	// repositories
	fcmJobRepo := repositories.NewFcmJobRepo(database)
	// services
	queueService := services.NewQueueService(queueConn, firebaseClient, fcmJobRepo)

	// worker
	queueWorker := workers.NewQueueService(queueService)
	go queueWorker.NotificationConsume()

	// controllers
	queueCtr := controllers.NewQueueService(queueService)

	// Restful API
	r := gin.Default()
	r.POST("/producer", queueCtr.SendMessage)
	panic(r.Run())
}
