package repositories

import (
	"github.com/tony-zhuo/consume-fcm/internal/models"
	"gorm.io/gorm"
)

type FcmJobRepo interface {
	Create(data models.FcmJob) error
}

type fcmJobRepo struct {
	db *gorm.DB
}

func NewFcmJobRepo(db *gorm.DB) FcmJobRepo {
	return &fcmJobRepo{
		db: db,
	}
}

func (r *fcmJobRepo) Create(data models.FcmJob) error {
	res := r.db.Create(data)
	return res.Error
}
