package eventmanager

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	"github.com/http"
)

func getURL() string {
	return fmt.Sprintf("%s://%s:%s", env.Get("EVENT_MANAGER_SCHEME", "http"), env.Get("EVENT_MANAGER_HOST", "localhost"), env.Get("EVENT_MANAGER_PORT", "8081"))
}

func Push(evt *Event) *exceptions.Exception {
	_, err := http.Put(fmt.Sprintf("%s/sendevent", getURL()), evt)
	if err != nil {
		return exceptions.NewComponentException(err)
	}
	return nil
}
