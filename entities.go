// Copyright (c) 2014 Jason Goecke
// entities.go

// To be implemented:
// 	- POST/Create EntityValueExpression
//	- DELETE EntityValueExpress

package wit

import (
	"encoding/json"
)

// Represents an Entity for the Wit API (https://wit.ai/docs/api#toc_15)
type Entity struct {
	Builtin bool   `json:"builtin"`
	Doc     string `json:"doc"`
	Id      string `json:"id"`
	Values  []EntityValue
}

// Represents a Value within an Entity
type EntityValue struct {
	Value       string   `json:"value"`
	Expressions []string `json:"expressions"`
}

// Represents a slice of entites when returend as an array (https://wit.ai/docs/api#toc_15)
type Entities []string

// Creates a new entity (https://wit.ai/docs/api#toc_19)
//
//		result, err := client.CreateEntity(entity)
func (client *WitClient) CreateEntity(entity *Entity) ([]byte, error) {
	data, err := json.Marshal(entity)
	result, statusCode, err := post(client.ApiBase+"/entities", data)
	if statusCode != 200 {
		return nil, err
	}
	return result, nil
}

// Deletes an entity (https://wit.ai/docs/api#toc_30)
//
//		result, err := client.DeleteEntity("favorite_city")
func (client *WitClient) DeleteEntity(id string) ([]byte, error) {
	result, statusCode, err := delete(client.ApiBase+"/entities/", id)
	if statusCode != 200 {
		return nil, err
	}
	return result, nil
}

// Deletes an entity's value (https://wit.ai/docs/api#toc_25)
//
// 		result, err := client.DeleteEntityValue("favorite_city", "Paris")
func (client *WitClient) DeleteEntityValue(id string, value string) ([]byte, error) {
	result, statusCode, err := delete(client.ApiBase+"/entities/", id+"/"+value)
	if statusCode != 200 {
		return nil, err
	}
	return result, nil
}

// Lists the configured entities (https://wit.ai/docs/api#toc_15)
//
//		result, err := client.Entities()
func (client *WitClient) Entities() (*Entities, error) {
	result, _, err := get(client.ApiBase + "/entities")
	if err != nil {
		return nil, err
	}
	entities, _ := parseEntities(result)
	return entities, nil
}

// Lists a single configured entity (https://wit.ai/docs/api#toc_17)
//
//		result, err := client.Entity("wit$temperature")
func (client *WitClient) Entity(id string) (*Entity, error) {
	result, _, err := get(client.ApiBase + "/entities/" + id)
	if err != nil {
		return nil, err
	}
	entity, _ := parseEntity(result)
	return entity, nil
}

// Updates and entity (https://wit.ai/docs/api#toc_22)
//
//		result, err := client.UpdateEntity(entity)
func (client *WitClient) UpdateEntity(entity *Entity) ([]byte, error) {
	data, err := json.Marshal(entity)
	result, statusCode, err := put(client.ApiBase+"/entities/"+entity.Id, data)
	if statusCode != 200 {
		return nil, err
	}
	return result, nil
}

// Parses the Entities JSON
func parseEntities(data []byte) (*Entities, error) {
	entities := &Entities{}
	err := json.Unmarshal(data, entities)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

// Parses the Entity JSON
func parseEntity(data []byte) (*Entity, error) {
	entity := &Entity{}
	err := json.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Parses the Entities Value JSON
func parseEntityValue(data []byte) (*EntityValue, error) {
	entityValue := &EntityValue{}
	err := json.Unmarshal(data, entityValue)
	if err != nil {
		return nil, err
	}
	return entityValue, nil
}
