// Copyright (c) 2014 Jason Goecke
// entities_test.go

package wit

import (
	"os"
	"testing"
)

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
		entity.Id != "wit$temperature" {
		t.Error("Message JSON did not parse properly.")
	}
}

func TestWitEntities(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	entities, err := client.Entities()
	if err != nil {
		t.Error("Did not fetch entities properly")
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
	}

	if entity.Id != "wit$age_of_person" ||
		entity.Builtin != true {
		t.Error("Did not parse entity properly")
	}
}

func TestCreateEntity(t *testing.T) {
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
	_, err = client.CreateEntity(entity)
	if err != nil {
		t.Error("Did not create entity properly")
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
	_, err = client.UpdateEntity(entity)
	if err != nil {
		t.Error("Did not parse entity properly")
	}
	if err != nil {
		t.Error("Did not update entity properly")
	}
}

func TestCreateEntityValue(t *testing.T) {
	data := `
	{
	  {
	    "value": "Paris",
	    "expressions": ["Paris", "City of Light", "Capital of France"]
	  }
	}`

	_, err := parseEntityValue([]byte(data))
	if err != nil {
		t.Error("Did not parse entity value properly")
	}
	// client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
}

func TestDeleteEntityValue(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	_, err := client.DeleteEntityValue("favorite_city", "Paris")
	if err != nil {
		t.Error("Did not delete entity value properly")
	}
}

func TestDeleteEntityValueExp(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	_, err := client.DeleteEntityValueExp("favorite_city", "Paris", "City of Light")
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteEntity(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	_, err := client.DeleteEntity("favorite_city")
	if err != nil {
		t.Error("Did not delete entity properly")
	}
}
