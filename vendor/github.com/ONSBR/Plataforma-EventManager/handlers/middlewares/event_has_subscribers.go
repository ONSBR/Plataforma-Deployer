package middlewares

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/actions"
	"github.com/ONSBR/Plataforma-EventManager/infra"
	"github.com/ONSBR/Plataforma-EventManager/processor"
)

//EventHasSubscribers checks if an event has at least one operation binding
func EventHasSubscribers(c *processor.Context) error {
	if refused, err := actions.RefuseEvent(c.Event); err != nil {
		return infra.NewSystemException(err.Error())
	} else if refused {
		return infra.NewSubscriberNotFoundException(fmt.Sprintf("Event %s has no subscribers", c.Event.Name))
	}
	return nil
}
