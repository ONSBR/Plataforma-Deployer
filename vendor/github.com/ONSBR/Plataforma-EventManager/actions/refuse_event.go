package actions

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/infra"
)

//RefuseEvent will refuse or not an event that did not have binding on platform
func RefuseEvent(event *domain.Event) (bool, error) {
	if event.IsSystemEvent() {
		return false, nil
	}
	if len(event.Bindings) == 0 {
		return true, infra.NewSubscriberNotFoundException(fmt.Sprintf("Event %s has no subscribers", event.Name))
	}
	return false, nil
}
