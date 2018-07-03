package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
)

//UpdateSplitState updates current split state and replace document on process memory
func UpdateSplitState(event *domain.Event, splitState *domain.SplitState, status string) error {
	for _, eventState := range splitState.Events {
		if eventState.EventOut == event.Name && eventState.Branch == event.Branch {
			eventState.Status = status
		}
	}
	return sdk.ReplaceDocument("splitted_events", map[string]string{"tag": event.Tag}, splitState)
}
