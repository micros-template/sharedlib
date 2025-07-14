package dto

type MailNotificationMessage struct {
	Receiver []string `json:"receiver"`
	MsgType  string   `json:"message_type"`
	Message  string   `json:"message"`
}
