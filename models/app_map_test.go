package models

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	yaml "gopkg.in/yaml.v2"
)

func TestShouldParseMapFile(t *testing.T) {
	Convey("should parse map file", t, func() {
		yml := `
Operacao:
    model: operacao
    fields:
        value:
            column: valor
        type:
            column: tipo
        personId:
            column: titular_id
        date:
            column: data
    filters:
        byTitular: "titular_id = :personId"

Conta:
    model: conta
    fields:
        personId:
            column: id
        balance:
            column: saldo
    filters:
        byPersonId: "id = :personId"`

		appMap := NewAppMap()
		err := yaml.Unmarshal([]byte(yml), &appMap)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		} else if len(appMap) == 0 {
			t.Fail()
		} else if _, ok := appMap["Conta"]; !ok {
			t.Fail()
		} else if _, ok := appMap["Operacao"]; !ok {
			t.Fail()
		} else {
			conta := appMap["Conta"]
			So(conta.Model, ShouldEqual, "conta")

			operacao := appMap["Operacao"]
			So(operacao.Model, ShouldEqual, "operacao")

			So(len(conta.Fields), ShouldBeGreaterThan, 0)

			So(len(operacao.Fields), ShouldEqual, 4)

			So(len(conta.Filters), ShouldEqual, 1)
		}

		fmt.Println(appMap)
	})

}
