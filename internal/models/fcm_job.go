package models

import (
	"time"
)

type FcmJob struct {
	Identifier string    `json:"identifier"`
	DeliverAt  time.Time `json:"deliver_at"`
}
