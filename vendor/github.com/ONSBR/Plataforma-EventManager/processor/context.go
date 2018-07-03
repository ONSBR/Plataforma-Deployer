package processor

import (
	"github.com/ONSBR/Plataforma-EventManager/bus"
	"github.com/ONSBR/Plataforma-EventManager/domain"
)

//Context is structure that have the current state
type Context struct {
	Event      *domain.Event
	dispatcher bus.Dispatcher
	Session    map[string]string
}

//Publish sends message to broker
func (c *Context) Publish(routingKey string, event *domain.Event) error {
	return c.dispatcher.Publish(routingKey, event.ToCeleryMessage())
}

//Dispatcher returns dispatcher from context
func (c *Context) Dispatcher() bus.Dispatcher {
	return c.dispatcher
}
