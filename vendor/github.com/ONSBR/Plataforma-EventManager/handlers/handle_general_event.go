package handlers

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/actions"
	"github.com/ONSBR/Plataforma-EventManager/processor"
	log "github.com/sirupsen/logrus"
)

//HandleGeneralEvent handle general event
func HandleGeneralEvent(c *processor.Context) error {
	log.Debug(fmt.Sprintf("HandleGeneralEvent %s on branch %s", c.Event.Name, c.Event.Branch))
	if c.Event.IsReproduction() {
		return c.Publish("store.executor", c.Event)
	}
	if events, err := actions.SplitEvent(c.Event); err != nil {
		return err
	} else if err := actions.SaveSplitState(events); err != nil {
		return err
	} else {
		for _, event := range events {
			if err := c.Publish("store.executor", event); err != nil {
				return err
			}
		}
	}
	return nil
}
