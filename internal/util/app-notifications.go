package util

type Notifications struct {
	MaxNotifications int
	Messages         []string
}

func NewNotifications() *Notifications {
	return &Notifications{
		MaxNotifications: 10,
		Messages:         make([]string, 0),
	}
}

func (n *Notifications) AddNotification(notification string) {
	n.Messages = append(n.Messages, notification)
	if len(n.Messages) > n.MaxNotifications {
		n.Messages = n.Messages[1:]
	}
}

var AppNotifications = NewNotifications()
