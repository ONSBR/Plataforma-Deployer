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

//FindBySystemID finds entity by SystemID on apicore
func FindBySystemID(entity, id string, response interface{}) *exceptions.Exception {
	filter := Filter{
		Entity: entity,
		Map:    "core",
		Name:   "bySystemId",
		Params: []Param{Param{
			Key:   "systemId",
			Value: id,
		},
		},
	}
	return Query(filter, response)
}
