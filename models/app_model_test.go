package models

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	yaml "gopkg.in/yaml.v2"
)

func TestShouldLoadAppModelFromYAML(t *testing.T) {
	Convey("should load app model from yaml", t, func() {
		yml := `
conta:
  titular:
    - string
  rg:
    - string
  saldo:
    - integer
`
		appmodel := make(AppModel)
		err := yaml.Unmarshal([]byte(yml), &appmodel)
		So(err, ShouldBeNil)
		So(appmodel["conta"], ShouldNotBeNil)
		So(appmodel["conta"]["titular"], ShouldNotBeNil)
		So(len(appmodel["conta"]["titular"]), ShouldEqual, 1)
	})
}
