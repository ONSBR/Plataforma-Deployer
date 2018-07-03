package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/infra"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
)

//SaveSplitState persist events executing state on process memory to track execution flow
func SaveSplitState(events []*domain.Event) error {
	splitState := domain.NewSplitState()
	if len(events) == 0 {
		return infra.NewArgumentException("cannot save split state document with no events")
	}
	splitState.Tag = events[0].Tag
	for _, event := range events {
		state := event.ToEventState()
		operation := event.Bindings[0]
		state.EventOut = operation.EventOut
		state.Version = operation.Version
		splitState.AddEventState(state)
	}
	return sdk.SaveDocument("splitted_events", splitState)
}
