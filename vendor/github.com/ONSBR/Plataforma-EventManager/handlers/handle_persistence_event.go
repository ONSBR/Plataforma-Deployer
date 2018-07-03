package handlers

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/bus"
	"github.com/ONSBR/Plataforma-EventManager/infra"
	"github.com/ONSBR/Plataforma-EventManager/processor"
	log "github.com/sirupsen/logrus"
)

//HandlePersistenceEvent handle persistence event to control execution flow
func HandlePersistenceEvent(c *processor.Context) error {
	log.Debug(fmt.Sprintf("HandlePersistenceEvent %s on branch %s", c.Event.Name, c.Event.Branch))
	_, err := c.Dispatcher().First(bus.EventPersistRequestQueue)
	if err != nil && err.Error() == infra.PersistEventQueueEmpty {
		return c.Publish("store.executor.inexecution", c.Event)
	} else if err != nil {
		log.Error(err.Error())
		return err
	}
	return c.Publish("persist.inexecution", c.Event)
}
