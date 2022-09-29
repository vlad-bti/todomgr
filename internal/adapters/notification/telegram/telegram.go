package telegram

import "testcode/test3/pkg/logger"

type telegramNotification struct {
	log *logger.Logger
}

func NewTelegramNotification(log *logger.Logger) *telegramNotification {
	return &telegramNotification{log}
}

func (r *telegramNotification) Send(msg interface{}) {
	r.log.Info("TelegramNotification - Send: %+v", msg)
}
