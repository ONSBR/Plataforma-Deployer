package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
)

//PushEventToEventStore from broker
func PushEventToEventStore(event *domain.Event) error {
	return SaveEventToStore(*event)
}
