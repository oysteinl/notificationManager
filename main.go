package notification

import (
	"github.com/oysteinl/notificationManager/notifier"
	"net/http"
	"time"
)

func NotifyTelegram(operation int8, messageToSend string, image []byte, config notifier.TelegramConfig) error {
	telegram := notifier.TelegramNotification{Method: operation, MessageToSend: messageToSend, Image: image, Config: config}
	return telegram.Send(&http.Client{Timeout: time.Second * 2})
}
