package webhook

import "testcode/test3/pkg/logger"

type webhookNotification struct {
	log *logger.Logger
}

func NewWebhookNotification(log *logger.Logger) *webhookNotification {
	return &webhookNotification{log}
}

func (r *webhookNotification) Send(msg any) {
	r.log.Info("WebhookNotification - Send: %+v", msg)
}
