// Copyright (c) 2014 Jason Goecke
// entities.go

package wit

import (
	"encoding/json"
	"net/url"
	"strings"
)

// Represents an Entity for the Wit API (https://wit.ai/docs/api#toc_15)
type Entity struct {
	Builtin bool          `json:"builtin"`
	Doc     string        `json:"doc"`
	Id      string        `json:"id"`
	Name	string 		  `json:"name"`
	Values  []EntityValue `json:"values"`
}

// Represents a Value within an Entity
type EntityValue struct {
	Value       string   `json:"value"`
	Expressions []string `json:"expressions"`
}

type Expression struct {
	Expression string `json:"expression"`
}

// Represents a slice of entites when returend as an array (https://wit.ai/docs/api#toc_15)
type Entities []string

// Creates a new entity (https://wit.ai/docs/api#toc_19)
//
//		result, err := client.CreateEntity(entity)
func (client *WitClient) CreateEntity(entity *Entity) (*Entity, error) {
	data, _ := json.Marshal(entity)
	result, err := post(client.ApiBase+"/entities", data)
	if err != nil {
		return nil, err
	}
	entity = &Entity{}
	err = json.Unmarshal(result, entity)
	return entity, nil
}

// Creates a new entity value (https://wit.ai/docs/api#toc_25)
//
//		result, err := client.CreateEntityValue("favorite_city, entityValue)
func (client *WitClient) CreateEntityValue(id string, entityValue *EntityValue) (*Entity, error) {
	data, _ := json.Marshal(entityValue)
	result, err := post(client.ApiBase+"/entities/"+id+"/values", data)
	if err != nil {
		return nil, err
	}
	entity := &Entity{}
	err = json.Unmarshal(result, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Creates a new entity value expression (https://wit.ai/docs/api#toc_25)
//
//		result, err := client.CreateEntityValueExp("favorite_city", "Barcelona", "Paella")
func (client *WitClient) CreateEntityValueExp(id string, value string, exp string) (*Entity, error) {
	jsonData, _ := json.Marshal(&Expression{exp})
	result, err := post(client.ApiBase+"/entities/"+id+"/values/"+value+"/expressions", jsonData)
	if err != nil {
		return nil, err
	}
	entity := &Entity{}
	err = json.Unmarshal(result, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Deletes an entity (https://wit.ai/docs/api#toc_30)
//
//		result, err := client.DeleteEntity("favorite_city")
func (client *WitClient) DeleteEntity(id string) error {
	id = url.QueryEscape(id)
	_, err := delete(client.ApiBase+"/entities", id)
	if err != nil {
		return err
	}
	return nil
}

// Deletes an entity's value (https://wit.ai/docs/api#toc_25)
//
// 		result, err := client.DeleteEntityValue("favorite_city", "Paris")
func (client *WitClient) DeleteEntityValue(id string, value string) ([]byte, error) {
	id = url.QueryEscape(id)
	result, err := delete(client.ApiBase+"/entities", id+"/values/"+value)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Deletes an entity's value's expression (https://wit.ai/docs/api#toc_35)
//
// 		result, err := client.DeleteEntityValueExp("favorite_city", "Paris", "")
func (client *WitClient) DeleteEntityValueExp(id string, value string, exp string) ([]byte, error) {
	id = url.QueryEscape(id)
	exp = strings.Replace(url.QueryEscape(exp), "+", "%20", -1)
	result, err := delete(client.ApiBase+"/entities", id+"/values/"+value+"/expressions/"+exp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Lists the configured entities (https://wit.ai/docs/api#toc_15)
//
//		result, err := client.Entities()
func (client *WitClient) Entities() (*Entities, error) {
	result, err := get(client.ApiBase + "/entities")
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
	id = url.QueryEscape(id)
	result, err := get(client.ApiBase + "/entities/" + id)
	if err != nil {
		return nil, err
	}
	entity, err := parseEntity(result)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Updates and entity (https://wit.ai/docs/api#toc_22)
//
//		result, err := client.UpdateEntity(entity)
func (client *WitClient) UpdateEntity(entity *Entity) ([]byte, error) {
	data, err := json.Marshal(entity)
	result, err := put(client.ApiBase+"/entities/"+entity.Id, data)
	if err != nil {
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
