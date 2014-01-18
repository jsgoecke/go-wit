// Copyright (c) 2014 Jason Goecke
// intents.go

package wit

import (
	"encoding/json"
)

// Represents intents in the Wit API (https://wit.ai/docs/api#toc_13)
type Intents []struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Doc      string `json:"doc"`
	Metadata string `json:"metadata"`
}

// Lists intents configured in the Wit API (https://wit.ai/docs/api#toc_13)
//
//		result, err := client.Intents()
func (client *WitClient) Intents() (*Intents, error) {
	result, _, err := get(client.ApiBase + "/intents")
	if err != nil {
		return nil, err
	}
	intents, _ := parseIntents(result)
	return intents, nil
}

// Parses the JSON for an Intent
func parseIntents(data []byte) (*Intents, error) {
	intents := &Intents{}
	err := json.Unmarshal(data, intents)
	if err != nil {
		return nil, err
	}
	return intents, nil
}
