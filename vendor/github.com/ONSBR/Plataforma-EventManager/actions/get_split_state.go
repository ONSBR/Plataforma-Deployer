package actions

import (
	"encoding/json"
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/infra"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
)

//GetSplitState retrieve split state based on event tag
func GetSplitState(event *domain.Event) (*domain.SplitState, error) {
	response, err := sdk.GetDocument("splitted_events", map[string]string{"tag": event.Tag})
	if err != nil {
		return nil, err
	}
	splitState := make([]domain.SplitState, 0)
	if err := json.Unmarshal([]byte(response), &splitState); err != nil {
		return nil, infra.NewArgumentException(err.Error())
	} else if len(splitState) == 0 {
		return nil, infra.NewSystemException(fmt.Sprintf("splitted_events with tag %s from event %s not found", event.Tag, event.Name))
	}
	return &splitState[0], nil
}
