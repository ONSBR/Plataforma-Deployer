package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
)

//SetFailureProcess on ApiCore
func SetFailureProcess(event *domain.Event) (err error) {
	id := event.Payload["instance_id"].(string)
	err = sdk.UpdateProcessInstance(id, "failure")
	return
}
