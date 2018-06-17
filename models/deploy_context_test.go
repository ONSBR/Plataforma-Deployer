package models

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	yaml "gopkg.in/yaml.v2"
)

func TestShouldUnMarshalMetadataFile(t *testing.T) {
	Convey("should parse metadata yaml file to AppMetadata model", t, func() {
		meta := NewAppMetadata()
		yml := `
operations:
  - name: consolida-saldo
    event: consolida.saldo.request
    commit: true`
		err := yaml.Unmarshal([]byte(yml), meta)
		if err != nil {
			fmt.Println(err.Error())
			t.Fail()
		} else if len(meta.Operations) == 0 {
			fmt.Println(meta)
			t.Fail()
		} else if meta.Operations[0].Name != "consolida-saldo" {
			fmt.Println(meta)
			t.Fail()
		}
	})

}
