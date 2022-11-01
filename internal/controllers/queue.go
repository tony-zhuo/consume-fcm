package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tony-zhuo/consume-fcm/internal/services"
	"net/http"
)

type QueueController interface {
	SendMessage(ctx *gin.Context)
}

type queueController struct {
	queueService services.QueueService
}

func NewQueueService(queueService services.QueueService) QueueController {
	return &queueController{
		queueService: queueService,
	}
}

func (c *queueController) SendMessage(ctx *gin.Context) {
	msg := services.MessageQueue{}
	fmt.Println("controller SendMessage start")
	if err := ctx.BindJSON(&msg); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	fmt.Println("controller SendMessage arg", msg)

	if err := c.queueService.Producer(msg); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
