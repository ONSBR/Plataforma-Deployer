package bus

import "github.com/ONSBR/Plataforma-EventManager/domain"

type DispatcherMock struct {
	OnPublish func(string, interface{}) error

	OnGet func(string, func(*domain.Event) error) error

	OnSwap func(string, string) error

	OnPop func(string) (*domain.Event, error)

	OnFirst func(string) (*domain.Event, error)
}

func (d *DispatcherMock) Publish(routingKey string, message interface{}) error {
	return d.OnPublish(routingKey, message)
}

func (d *DispatcherMock) Get(queue string, action func(*domain.Event) error) error {
	return d.OnGet(queue, action)
}

func (d *DispatcherMock) Swap(queueFrom string, routingKey string) error {
	return d.OnSwap(queueFrom, routingKey)
}

func (d *DispatcherMock) Pop(queue string) (*domain.Event, error) {
	return d.OnPop(queue)
}
func (d *DispatcherMock) First(queue string) (*domain.Event, error) {
	return d.OnFirst(queue)
}
