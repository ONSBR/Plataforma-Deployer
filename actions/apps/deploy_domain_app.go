package apps

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ONSBR/Plataforma-Deployer/git"
	"github.com/ONSBR/Plataforma-Deployer/models"
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
	yaml "gopkg.in/yaml.v2"
)

type templateData struct {
	DatabaseName string
	Models       models.AppModel
}

func deployDomainAppWorker(queue chan *models.DeployContext) {
	for context := range queue {
		context.Start(doDomainDeploy)
	}
}

func doDomainDeploy(context *models.DeployContext) *exceptions.Exception {
	context.RemoveContainer(context.GetContainerName())
	if ex := downloadTemplate(context); ex != nil {
		return ex
	} else if entities, ex := loadDomainEntities(context); ex != nil {
		return ex
	} else if ex := compileApp(context, entities); ex != nil {
		return ex
	} else if ex := saveDomainPersistOperation(context); ex != nil {
		return ex
	}
	return nil
}

func saveDomainPersistOperation(context *models.DeployContext) *exceptions.Exception {
	return context.PersistOperations([]*models.Operation{
		&models.Operation{
			Event:  fmt.Sprintf("%s.persist.request", context.Info.App.SystemID),
			Name:   fmt.Sprintf("%s.persist", context.Info.App.Name),
			Commit: true,
		},
		&models.Operation{
			Event:  fmt.Sprintf("%s.merge.request", context.Info.App.SystemID),
			Name:   fmt.Sprintf("%s.merge", context.Info.App.Name),
			Commit: true,
		},
	})
}

func compileApp(context *models.DeployContext, entities []models.AppModel) *exceptions.Exception {
	if compiled, err := compile(context, entities); err != nil {
		return err
	} else {
		path := fmt.Sprintf("%s/domain.py", getModelPath(context))
		fd, err := os.Open(path)
		if err == nil {
			os.Remove(path)
		} else {
			fd.Close()
		}
		if err := ioutil.WriteFile(path, []byte(compiled), 0666); err != nil {
			return exceptions.NewComponentException(err)
		}
		//redirect deployer to point to compiled app instead of domain app
		context.RootPath = getAppPath(context)
	}
	return nil
}

func getModelPath(context *models.DeployContext) string {
	return fmt.Sprintf("%s/Plataforma-Domain/Platform.App/python-template/model", getTemplatePath(context))
}

func getAppPath(context *models.DeployContext) string {
	return fmt.Sprintf("%s/Plataforma-Domain/Platform.App/python-template", getTemplatePath(context))
}

func compile(context *models.DeployContext, list []models.AppModel) (string, *exceptions.Exception) {
	template := fmt.Sprintf("%s/domain_remote.tmpl", getModelPath(context))
	data, err := ioutil.ReadFile(template)
	if err != nil {
		return "", exceptions.NewComponentException(err)
	}
	tmpl := templateData{
		DatabaseName: strings.ToLower(context.Info.App.SystemName),
	}
	tmpl.Models = make(models.AppModel)
	for _, mod := range list {
		for k := range mod {
			tmpl.Models[k] = mod[k]
		}
	}
	return applyTemplate(string(data), tmpl)
}

func applyTemplate(appTemplate string, data templateData) (string, *exceptions.Exception) {
	buf := bytes.NewBufferString("")
	tmpl := template.Must(template.New(appTemplate).Funcs(templateHelpers).Parse(appTemplate))
	err := tmpl.Execute(buf, data)
	if err != nil {
		return "", exceptions.NewComponentException(err)
	}
	return buf.String(), nil
}

func loadDomainEntities(context *models.DeployContext) ([]models.AppModel, *exceptions.Exception) {
	path := fmt.Sprintf("%s/Dominio", context.GetWorkspace())
	if files, err := ioutil.ReadDir(path); err != nil {
		return nil, exceptions.NewComponentException(err)
	} else if len(files) == 0 {
		return nil, exceptions.NewInvalidArgumentException(fmt.Errorf("No entity description found in app"))
	} else {
		list := make([]models.AppModel, 0)
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			if data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, f.Name())); err != nil {
				return nil, exceptions.NewComponentException(err)
			} else {
				appmodel := make(models.AppModel)
				if err := yaml.Unmarshal(data, &appmodel); err != nil {
					return nil, exceptions.NewInvalidArgumentException(err)
				}
				list = append(list, appmodel)
			}
		}
		return list, nil
	}
}

func getTemplatePath(context *models.DeployContext) string {
	return fmt.Sprintf("/worker/deploys/%s", context.Info.ID)
}

func downloadTemplate(context *models.DeployContext) *exceptions.Exception {
	path := getTemplatePath(context)
	fd, err := os.Open(path)
	if err == nil {
		fd.Close()
		os.RemoveAll(path)
	}
	if ex := git.CloneRepo(getTemplatePath(context), "https://github.com/ONSBR/Plataforma-Domain", "feature/deploy-remoto"); ex != nil {
		return ex
	}
	return nil
}

var templateHelpers = template.FuncMap{
	"join": func(l []string, sep string) string {
		return strings.Join(l, sep)
	},
	"keys": func(m models.AppModelEntity) []string {
		l := make([]string, len(m))
		i := 0
		for k := range m {
			l[i] = k
			i++
		}
		return l
	},
	"ormType": func(params []string) (string, error) {
		if len(params) == 0 {
			return "", nil
		}
		dataType := params[0]
		t, ok := ormTypes[dataType]
		if !ok {
			return "", fmt.Errorf("Type %s is not defined", dataType)
		}
		return t, nil
	},
}

var ormTypes = map[string]string{
	"string":    "String",
	"integer":   "Integer",
	"char":      "Char",
	"text":      "Text",
	"bigint":    "BigInt",
	"float":     "Float",
	"real":      "Real",
	"double":    "Float",
	"decimal":   "Decimal",
	"boolean":   "Boolean",
	"time":      "Time",
	"date":      "Date",
	"datetime":  "DateTime",
	"hstore":    "HsStore",
	"json":      "Json",
	"jsonb":     "Jsonb",
	"blob":      "Blob",
	"uuid":      "sap.UUID(as_uuid=True)",
	"uuidV1":    "UUIDV1",
	"uuidV4":    "UUIDV4",
	"timestamp": "TIMESTAMP",
}
