package eventstore

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
)

func Push(event domain.Event) error {
	return influxPush(event)
}

func Count(field, value, last string) int {
	return totalEventsByField(field, value, last)
}

func Query(field, value, last string) []domain.Event {
	return queryEvents(field, value, last)
}
