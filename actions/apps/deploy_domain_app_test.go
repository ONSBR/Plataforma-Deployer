package apps

import (
	"fmt"
	"testing"

	"github.com/ONSBR/Plataforma-Deployer/models"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldCompileTemplate(t *testing.T) {
	Convey("should compile template with app entities", t, func() {
		model := make(models.AppModel)
		model["conta"] = make(models.AppModelEntity)
		model["conta"]["titular"] = make([]string, 1)
		model["conta"]["titular"][0] = "string"
		model["conta"]["saldo"] = make([]string, 1)
		model["conta"]["saldo"][0] = "integer"
		model["conta"]["titular_id"] = make([]string, 1)
		model["conta"]["titular_id"][0] = "uuid"

		model["pessoa"] = make(models.AppModelEntity)
		model["pessoa"]["nome"] = make([]string, 1)
		model["pessoa"]["nome"][0] = "string"
		model["pessoa"]["rg"] = make([]string, 1)
		model["pessoa"]["rg"][0] = "string"

		data := templateData{
			DatabaseName: "bankapp",
			Models:       model,
		}
		appTemplate := `
from database import Base
from uuid import uuid4
from core.temporal.models import TemporalModelMixin
import sqlalchemy.dialects.postgresql as sap
from sqlalchemy.ext.declarative import declared_attr
from sqlalchemy import *
from datetime import datetime

def get_db_name():
    return "{{.DatabaseName}}"

{{range $model, $attrs := .Models}}
class {{$model}}(Base, TemporalModelMixin):

    def __init__(self, rid=None, id=None, deleted=False, meta_instance_id=None, {{range $attr, $params := $attrs}}{{$attr}}=None,{{end}} _metadata=None, **kwargs):
        self.rid = rid
        self.id = id
        self.deleted = deleted
        self.meta_instance_id = meta_instance_id
        {{range $attr, $params := $attrs}}
        self.{{$attr}} = {{$attr}}
        {{end}}
        self._metadata = _metadata
        self.branch = kwargs.get('branch', 'master')
        self.from_id = kwargs.get('from_id')
        self.modified = kwargs.get('modified')

    def dict(self):
        return { {{range $attr, $params := $attrs}}
            "{{$attr}}":self.{{$attr}},{{end}}
            "id": self.id,
            "branch":self.branch,
            "modified":self.modified,
            "created_at":self.created_at,
            "_metadata": self._metadata
        }
    @declared_attr
    def __tablename__(cls):
        return cls.__name__.lower()

    class Temporal:
        fields = ('deleted','modified','created_at', 'meta_instance_id', 'from_id', 'branch', {{range $attr, $params := $attrs}} '{{$attr}}', {{end}} )


{{range $attr, $params := $attrs }}
    {{$attr}} = Column({{ormType $params}}){{end}}
    id = Column(sap.UUID(as_uuid=True), default=uuid4)
    deleted = Column(sap.BOOLEAN(), default=False)
    meta_instance_id = Column(sap.UUID(as_uuid=True))
    modified = Column(DateTime(), default=datetime.utcnow())
    created_at = Column(DateTime(), default=datetime.utcnow())
    branch = Column(String(), default='master')
    from_id = Column(sap.UUID(as_uuid=True), nullable=True)
{{end}}
	`
		_, err := applyTemplate(appTemplate, data)
		if err != nil {
			fmt.Println(err.Error())
			t.Fail()
		}
	})
}
