package producer

import "lms-be/event/model"

// Producer represents an event producer interface.
type Producer interface {
	Publish(request model.PublishRequest)
}
