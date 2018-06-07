package apicore

import (
	"github.com/ONSBR/Plataforma-Deployer/models/exceptions"
)

//FindByID finds entity by ID on apicore
func FindByID(entity, id string, response interface{}) *exceptions.Exception {
	filter := Filter{
		Entity: entity,
		Map:    "core",
		Name:   "byId",
		Params: []Param{Param{
			Key:   "id",
			Value: id,
		},
		},
	}
	return Query(filter, response)
}
