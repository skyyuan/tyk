package main

import (
	"encoding/json"
)

type NotificationCommand string

const (
	NoticeApiUpdated    NotificationCommand = "ApiUpdated"
	NoticeApiRemoved    NotificationCommand = "ApiRemoved"
	NoticeApiAdded      NotificationCommand = "ApiAdded"
	NoticeGroupReload   NotificationCommand = "GroupReload"
	NoticePolicyChanged NotificationCommand = "PolicyChanged"
)

// Notification is a type that encodes a message published to a pub sub channel
type Notification struct {
	Command NotificationCommand `json:"command"`
	Payload string              `json:"payload"`
}

// Notifier is an interface that sends notifications
type Notifier interface {
	Notify(string) bool
}

// RedisNotifier implements Notifier and will use redis pub/sub channles to send notifications
type RedisNotifier struct {
	store   *RedisClusterStorageManager
	channel string
}

// Notify will send a notification to a channel
func (r *RedisNotifier) Notify(notification Notification) bool {
	toSend, err := json.Marshal(notification)
	if err != nil {
		log.Error("Problem marshalling notification!")
		log.Error(err)
		return false
	} else {
		log.Debug("Sending notification", notification)
		sentErr := r.store.Publish(r.channel, string(toSend))
		if sentErr != nil {
			log.Error("Could not send notification")
			log.Error(err)
			return false
		}
	}
	return true
}
