package eventmanager

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/PMoneda/http"
)

func getURL() string {
	return fmt.Sprintf("%s://%s:%s", env.Get("EVENT_MANAGER_SCHEME", "http"), env.Get("EVENT_MANAGER_HOST", "localhost"), env.Get("EVENT_MANAGER_PORT", "8081"))
}

//Push event to event manager
func Push(evt *domain.Event) error {
	_, err := http.Put(fmt.Sprintf("%s/sendevent", getURL()), evt)
	if err != nil {
		return err
	}
	return nil
}

//LastEvents returns last events from eventstore
func LastEvents(field, value, last string) ([]*domain.Event, error) {
	type Response struct {
		Result []*domain.Event `json:"result"`
	}
	r := Response{
		Result: make([]*domain.Event, 0),
	}
	err := http.GetJSON(fmt.Sprintf("%s/events?field=%s&value=%s&last=%s", getURL(), field, value, last), &r)
	if err != nil {
		return nil, err
	}
	return r.Result, nil
}
