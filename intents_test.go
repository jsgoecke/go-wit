// Copyright (c) 2014 Jason Goecke
// intents_test.go

package wit

import (
	//"os"
	"testing"
)

func TestWitIntentsParsing(t *testing.T) {
	data := `
	[ {
	  "id" : "1234",
	  "name" : "recover_password",
	  "doc" : "Recover password (which is different from Reset password).",
	  "metadata" : "password_23433253254"
	}, {
	  "id" : "52bab833-9a1e-4bff-b659-99ee95e6c1f9",
	  "name" : "transfer",
	  "doc" : "Transfer some amount of money between two accounts."
	}, {
	  "id" : "52bab833-3e23-4c67-9cfc-a0fed605bd77",
	  "name" : "show_movie",
	  "doc" : "Show a given movie."
	} ]`

	intents, err := parseIntents([]byte(data))
	if err != nil {
		t.Error(err.Error())
	}

	for cnt, intent := range *intents {
		switch cnt {
		case 0:
			if intent.ID != "1234" {
				t.Error("Intents JSON did not parse properly.")
			}
		case 1:
			if intent.Name != "transfer" {
				t.Error("Intents JSON did not parse properly.")
			}
		case 2:
			if intent.Doc != "Show a given movie." {
				t.Error("Intents JSON did not parse properly.")
			}
		}
	}
}

// temporary: the "good_bye" intent is not added to new instances anymore.
// Wit.AI will soon add a POST /intent endpoint to add new intents, uncomment this test then.

// func TestWitIntents(t *testing.T) {
// 	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
// 	intents, err := client.Intents()
// 	if err != nil {
// 		t.Error("Did not fetch intents properly")
// 	}

// 	goodBye := false
// 	for _, value := range *intents {
// 		if value.Name == "good_bye" {
// 			goodBye = true
// 		}
// 	}
// 	if goodBye != true {
// 		t.Error("Intents returned not expected")
// 	}
// }
