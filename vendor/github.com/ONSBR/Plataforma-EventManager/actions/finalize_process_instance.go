package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
)

//FinalizeProcessInstance on ApiCore
func FinalizeProcessInstance(event *domain.Event) (err error) {
	id := event.Payload["instance_id"].(string)
	err = sdk.UpdateProcessInstance(id, "finished")
	return
}
