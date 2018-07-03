package eventstore

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/domain"

	"github.com/ONSBR/Plataforma-EventManager/infra"
)

type Serie struct {
	Name    string          `json:"name"`
	Columns []string        `json:"columns"`
	Values  [][]interface{} `json:"values"`
}

type StmtResult struct {
	ID     int     `json:"statement_id"`
	Series []Serie `json:"series"`
}
type QueryResult struct {
	Results []StmtResult `json:"results"`
}

func queryCount(last string) string {
	return fmt.Sprintf(`SELECT count("count") AS "total"
	FROM "%s"."%s"."events"
	WHERE
	time > now() - %s`, infra.GetEnv("DATABASE", "teste"), infra.GetEnv("RETENTION_POLICY", "teste"), last)
}

func basicQuery(last string) string {
	return fmt.Sprintf(`SELECT "data"
	FROM "%s"."%s"."events"
	WHERE
	time > now() - %s`, infra.GetEnv("DATABASE", "teste"), infra.GetEnv("RETENTION_POLICY", "teste"), last)
}

func queryEvents(field, value, last string) []domain.Event {
	if last == "" {
		last = "1h"
	}
	q := fmt.Sprintf(`%s AND "%s"='%s'`, basicQuery(last), field, value)
	if field == "" {
		q = basicQuery(last)
	}
	s, _ := executeStatement(q)
	r := QueryResult{}
	if err := json.Unmarshal([]byte(s), &r); err != nil {
		return []domain.Event{}
	}
	if len(r.Results) > 0 && len(r.Results[0].Series) > 0 {
		c := r.Results[0].Series[0].Values
		result := make([]domain.Event, len(c))
		i := 0
		for _, reg := range c {
			switch v := reg[1].(type) {
			case string:
				decoded, _ := base64.StdEncoding.DecodeString(v)
				evt := new(domain.Event)
				json.Unmarshal(decoded, evt)
				evt.Timestamp = reg[0].(string)
				result[i] = *evt
			}
			i++
		}
		return result
	}
	return []domain.Event{}
}

func totalEventsByField(field, value, last string) int {
	q := fmt.Sprintf(`%s AND "%s"='%s'`, queryCount(last), field, value)
	s, _ := executeStatement(q)
	r := QueryResult{}
	if err := json.Unmarshal([]byte(s), &r); err != nil {
		return 0
	}
	if len(r.Results) > 0 && len(r.Results[0].Series) > 0 {
		c := r.Results[0].Series[0].Values[0][1]
		switch v := c.(type) {
		case float64:
			return int(v)
		default:
			return 0
		}
	}
	return 0

}
