package main

import (
	"fmt"
	"os"

	"github.com/ONSBR/Plataforma-EventManager/infra/factories"

	"github.com/ONSBR/Plataforma-EventManager/actions"
	"github.com/ONSBR/Plataforma-EventManager/api"
	"github.com/ONSBR/Plataforma-EventManager/bus"
	"github.com/ONSBR/Plataforma-EventManager/eventstore"
	log "github.com/sirupsen/logrus"
)

func init() {
	os.Setenv("DATABASE", "event_manager")
	os.Setenv("RETENTION_POLICY", "platform_events")
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func registerActionsToRabbitMq() *bus.Broker {
	broker := factories.GetBroker()
	broker.RegisterWorker(1, bus.EventstoreQueue, actions.PushEventToEventStore)
	broker.RegisterWorker(1, bus.EventProcessFinishedQueue, actions.FinalizeProcessInstance)
	broker.RegisterWorker(1, bus.EventExceptionQueue, actions.SetFailureProcess)
	broker.Listen()
	return broker
}

func main() {
	fmt.Println(logo())
	log.Info("Starting Event Manager")
	log.Info("Installing Bus")
	registerActionsToRabbitMq()
	log.Info("Starting Mongo")
	eventstore.Install()
	log.Info("Starting API")
	api.Build()
}

func logo() (logo string) {
	logo = `
 	 ______               _     __  __
	|  ____|             | |   |  \/  |
	| |____   _____ _ __ | |_  | \  / | __ _ _ __   __ _  __ _  ___ _ __
	|  __\ \ / / _ \ '_ \| __| | |\/| |/ _' | '_ \ / _' |/ _' |/ _ \ '__|
	| |___\ V /  __/ | | | |_  | |  | | (_| | | | | (_| | (_| |  __/ |
	|______\_/ \___|_| |_|\__| |_|  |_|\__,_|_| |_|\__,_|\__, |\___|_|
	                                                      __/ |
	                                                     |___/
	`
	return
}
