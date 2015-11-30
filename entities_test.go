// Copyright (c) 2014 Jason Goecke
// entities_test.go

package wit

import (
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"
)

var entityName string

func findEntityValue(e *Entity, v string) (value EntityValue, b bool) {
	for _, val := range e.Values {
		if val.Value == v {
			return val, true
		}
	}
	return value, false
}

func findStringInArray(a []string, v string) (s string, b bool) {
	for _, val := range a {
		if val == v {
			return val, true
		}
	}
	return s, false
}

func TestWitEntitiesParsing(t *testing.T) {
	data := `
	[
	   "wit$amount_of_money",
	   "wit$contact",
	   "wit$datetime",
	   "wit$on_off",
	   "wit$phrase_to_translate",
	   "wit$temperature"
	]`

	entities, err := parseEntities([]byte(data))
	if err != nil {
		t.Error(err.Error())
	}

	for cnt, value := range *entities {
		switch cnt {
		case 0:
			if value != "wit$amount_of_money" {
				t.Error("Entities JSON did not parse properly.")
			}
		case 3:
			if value != "wit$on_off" {
				t.Error("Entities JSON did not parse properly.")
			}
		case 5:
			if value != "wit$temperature" {
				t.Error("Entities JSON did not parse properly.")
			}
		}
	}
}

func TestWitEntityParsing(t *testing.T) {
	data := `
	{
	  "builtin": true,
	  "doc": "Temperature in degrees Celcius or Fahrenheit",
	  "id": "wit$temperature"
	}`

	entity, err := parseEntity([]byte(data))
	if err != nil {
		t.Error(err.Error())
	}

	if entity.Builtin != true ||
		entity.Doc != "Temperature in degrees Celcius or Fahrenheit" ||
		entity.ID != "wit$temperature" {
		t.Error("Message JSON did not parse properly.")
	}
}

func TestWitEntities(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	entities, err := client.Entities()
	if err != nil {
		t.Error(err)
		return
	}

	ageOfPerson := false
	for _, value := range *entities {
		if value == "wit$age_of_person" {
			ageOfPerson = true
		}
	}
	if ageOfPerson != true {
		t.Error("Entities returned not expected")
	}
}

func TestWitEntity(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	entity, err := client.Entity("wit$age_of_person")
	if err != nil {
		t.Error(err)
		return
	}
	if entity.Name != "age_of_person" ||
		entity.Builtin != true {
		t.Error("Did not parse entity properly")
	}

	// Now test for when the entity is not present
	_, err = client.Entity("age_of_person")
	if err.Error() != http.StatusText(404) {
		t.Error("Should have returned a not found error")
	}

}

func TestCreateEntity(t *testing.T) {
	// Create random string for testing
	rand.Seed(time.Now().UTC().UnixNano())
	alpha := "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	buf := make([]byte, 10)
	for i := 0; i < 10; i++ {
		buf[i] = alpha[rand.Intn(len(alpha))]
	}
	entityName = "favorite_city" + string(buf)

	data := `
	{
	  "doc": "A city that I like",
	  "id": "favorite_city",
	  "values": [
	    {
	      "value": "Paris",
	      "expressions": ["Paris", "City of Light", "Capital of France"]
	    }
	  ]
	}`

	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	entity, err := parseEntity([]byte(data))
	if err != nil {
		t.Error("Did not parse entity properly")
	}
	entity.ID = entityName
	entity, err = client.CreateEntity(entity)
	if err != nil {
		t.Errorf("Error creating entity %s", err.Error())
		return
	}
	if entity.Doc != "A city that I like" {
		t.Error("Entity was not created properly, doc not set")
	}
	if entity.Values[0].Value != "Paris" {
		t.Error("Entity was not created properly, values not set")
	}
	_, err = client.CreateEntity(entity)
	if err.Error() != http.StatusText(409) {
		t.Error("Expected a 409 since the entity already exists")
	}
}

func TestUpdateEntity(t *testing.T) {
	data := `
	{
	  "id": "favorite_city",
	  "doc": "These are cities worth going to",
	  "values": [
	    {
	      "value": "Paris",
	      "expressions": ["Paris", "City of Light", "Capital of France"]
	    },
	    {
	      "value": "Seoul",
	      "expressions": ["Seoul", "서울", "Kimchi paradise"]
	    }
	  ]
	}`

	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	entity, err := parseEntity([]byte(data))
	entity.ID = entityName
	_, err = client.UpdateEntity(entity)
	if err != nil {
		t.Error("Did not update entity properly")
	}
}

func TestCreateEntityValue(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	entityValue := &EntityValue{}
	entityValue.Value = "Barcelona"
	entityValue.Expressions = []string{"Med", "Sagrada Familia", "Gaudi"}
	entity, err := client.CreateEntityValue(entityName, entityValue)
	if err != nil {
		t.Error(err)
		return
	}
	barcelona, found := findEntityValue(entity, "Barcelona")
	if !found || barcelona.Value != "Barcelona" {
		t.Error("Did not add Barcelona to entity's value properly")
	}
	_, found = findStringInArray(barcelona.Expressions, "Sagrada Familia")
	if !found {
		t.Error("Did not add Sagrada Familia to entity's value expression properly")
	}
}

func TestCreateEntityValueExp(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	entity, err := client.CreateEntityValueExp(entityName, "Barcelona", "Paella")
	if err != nil {
		t.Error(err)
		return
	}
	barcelona, found := findEntityValue(entity, "Barcelona")
	if !found || barcelona.Value != "Barcelona" {
		t.Error("Did not add Barcelona to entity's value properly")
	}
	_, found = findStringInArray(barcelona.Expressions, "Paella")
	if !found {
		t.Error("Did not add Sagrada Familia to entity's value expression properly")
	}
}

func TestDeleteEntityValueExp(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	_, err := client.DeleteEntityValueExp(entityName, "Paris", "City of Light")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDeleteEntityValue(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	_, err := client.DeleteEntityValue(entityName, "Paris")
	if err != nil {
		t.Error("Did not delete entity value properly")
	}
}

func TestDeleteEntity(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	err := client.DeleteEntity("favorite_city")
	if err == nil {
		t.FailNow()
	}
	if err.Error() != http.StatusText(404) {
		t.Error("Delete should have returned 'Entity not found'")
	}
	err = client.DeleteEntity(entityName)
	if err != nil {
		t.Error(err)
		return
	}
}
