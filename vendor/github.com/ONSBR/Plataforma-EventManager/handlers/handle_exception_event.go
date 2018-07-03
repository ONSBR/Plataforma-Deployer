package handlers

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/actions"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/processor"
	log "github.com/sirupsen/logrus"
)

//HandleExceptionEvent handle exception events
func HandleExceptionEvent(c *processor.Context) error {
	log.Debug(fmt.Sprintf("HandleExceptionEvent %s on branch %s", c.Event.Name, c.Event.Branch))
	err := actions.SwapPersistEventToExecutorQueue(c.Dispatcher())
	splitState, err := actions.GetSplitState(c.Event)
	if err != nil {
		return err
	}
	if err := actions.UpdateSplitState(c.Event, splitState, domain.Error); err != nil {
		return err
	}
	return c.Publish("store.executor.exception", c.Event)
}
