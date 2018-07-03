package bus

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
)

//Dispatcher is a basic interface to send messages to broker
type Dispatcher interface {
	Publish(routingKey string, message interface{}) error

	Get(queue string, action func(*domain.Event) error) error

	Swap(queueFrom string, routingKey string) error

	Pop(queue string) (*domain.Event, error)

	First(queue string) (*domain.Event, error)
}
