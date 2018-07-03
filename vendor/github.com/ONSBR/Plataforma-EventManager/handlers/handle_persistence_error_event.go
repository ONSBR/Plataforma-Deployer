package handlers

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/actions"
	"github.com/ONSBR/Plataforma-EventManager/processor"
	log "github.com/sirupsen/logrus"
)

//HandlePersistenceErrorEvent handle persist error events to avoid blocking platform
func HandlePersistenceErrorEvent(c *processor.Context) error {
	log.Debug(fmt.Sprintf("HanlePersistenceErrorEvent %s on branch %s", c.Event.Name, c.Event.Branch))
	if err := actions.SwapPersistEventToExecutorQueue(c.Dispatcher()); err != nil {
		return err
	}
	return c.Publish("store.executor", c.Event)
}
