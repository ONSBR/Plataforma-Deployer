package eventstore

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/infra"
	"github.com/PMoneda/http"
	log "github.com/sirupsen/logrus"
)

//Install prepare influxdb by creating database and retenrion policy
func Install() {
	log.Info("Connecting to InfluxDB")
	go (func() {
		log.Info("Trying to connect to InfluxDB")
		for shouldCreateDatabase(infra.GetEnv("DATABASE", "teste")) {
			log.Info("Trying to create database")
			createDatabase(infra.GetEnv("DATABASE", "teste"))
			log.Info("Trying to create Retention Policy")
			createRetentionPolicy(infra.GetEnv("RETENTION_POLICY", "teste"), infra.GetEnv("DATABASE", "teste"))
			time.Sleep(5 * time.Second)
		}
		log.Info("Connected to influx")
	})()

}

//Push event to Influx
func influxPush(event domain.Event) error {
	if s, err := influxWrite(compileEventToLintProtocol(event)); err != nil {
		return err
	} else if s != "" {
		return errors.New(s)
	}
	return nil
}

func compileEventToLintProtocol(event domain.Event) string {
	//cpu_load_short,host=server02,region=us-west value=0.55 1422568543702900257
	s, _ := json.Marshal(event)
	parts := strings.Split(event.Name, ".")
	eventClass := parts[len(parts)-1]
	eventSubclass := "no_subclass"
	if len(parts) >= 2 {
		eventSubclass = parts[len(parts)-2]
	}
	encoded := base64.StdEncoding.EncodeToString(s)
	if event.InstanceID == "" {
		event.InstanceID = "new_instance"
	}
	cmd := fmt.Sprintf(`events,name=%s,scope=%s,tag=%s,branch=%s,instanceId=%s,class=%s,subclass=%s,owner=%s,appOrigin=%s count=1,data="%s"`, event.Name, event.Scope, event.Tag, event.Branch, event.InstanceID, eventClass, eventSubclass, event.Owner, event.AppOrigin, string(encoded))
	return cmd
}

func getBaseUrl() string {
	return fmt.Sprintf("http://%s:%s", infra.GetEnv("INFLUX_HOST", "localhost"), infra.GetEnv("INFLUX_PORT", "8086"))
}

func createDatabase(name string) {
	executeStatement(fmt.Sprintf("CREATE DATABASE %s", name))
}

func createRetentionPolicy(name, db string) {
	executeStatement(fmt.Sprintf("CREATE RETENTION POLICY %s ON %s DURATION 500w REPLICATION 1", name, db))
}

func shouldCreateDatabase(name string) bool {
	dbs := showDatabases()
	for _, db := range dbs {
		if db == name {
			return false
		}
	}
	return true
}

func showRetentionPolicy(dbname string) []string {
	r, err := executeStatement(fmt.Sprintf("SHOW RETENTION POLICIES ON %s", dbname))
	if err != nil {
		return []string{}
	}
	jsonParsed, _ := gabs.ParseJSON([]byte(r))
	b := jsonParsed.Path("results.series.values").String()
	//TODO refactor
	b = strings.Replace(b, "[", "", -1)
	b = strings.Replace(b, "]", "", -1)
	b = strings.Replace(b, "\"", "", -1)
	return strings.Split(b, ",")
}

func showDatabases() []string {
	r, err := executeStatement("SHOW DATABASES")
	if err != nil {
		return []string{}
	}
	jsonParsed, _ := gabs.ParseJSON([]byte(r))
	b := jsonParsed.Path("results.series.values").String()
	//TODO refactor
	b = strings.Replace(b, "[", "", -1)
	b = strings.Replace(b, "]", "", -1)
	b = strings.Replace(b, "\"", "", -1)
	return strings.Split(b, ",")
}

func executeStatement(stmt string) (string, error) {
	_url := fmt.Sprintf("%s/query?u=%s&p=%s", getBaseUrl(), infra.GetEnv("INFLUX_USER", ""), infra.GetEnv("INFLUX_PASSWORD", ""))
	payload := url.PathEscape(fmt.Sprintf("q=%s", stmt))
	return http.Post(_url, payload, http.Header{Key: "Content-Type", Value: "application/x-www-form-urlencoded"})

}

func influxWrite(point string) (string, error) {
	_url := fmt.Sprintf("%s/write?u=%s&p=%s&db=%s&rp=%s", getBaseUrl(), infra.GetEnv("INFLUX_USER", ""), infra.GetEnv("INFLUX_PASSWORD", ""), infra.GetEnv("DATABASE", "teste"), infra.GetEnv("RETENTION_POLICY", "teste"))
	return http.Post(_url, point)

}
