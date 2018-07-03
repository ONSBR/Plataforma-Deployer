package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/bus"
	log "github.com/sirupsen/logrus"
)

/*
Persist events are controled by a buffer queue so when a persist event comes
it will be store in two queues the first one: EventPersistQueue that enqueue all requesting
persist events when has some persist event running on platform
second one: EventPersistRequestQueue is control queue that have all persist queue running or awaiting for execution
When Event Manager receive a (done) event from persistence event we will swap event from awaiting queue and pop event from
pending queue
*/

//SwapPersistEventToExecutorQueue swaps event from event persist queue to executor queue
func SwapPersistEventToExecutorQueue(dispatcher bus.Dispatcher) error {
	log.Debug("Swapping persist event to executor queue")
	err := dispatcher.Swap(bus.EventPersistQueue, "executor.store")
	if err != nil {
		return err
	}
	_, err = dispatcher.Pop(bus.EventPersistRequestQueue)
	return err
}
