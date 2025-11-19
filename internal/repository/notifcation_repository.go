package repository

import (
	"context"
)

type NotificationRepository interface {
}

type notificationRepository struct {
}

func NewNotificationRepository() NotificationRepository {
	return &notificationRepository{}
}

func (n *notificationRepository) NewNotification(ctx context.Context) {

}
