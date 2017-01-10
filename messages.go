// Copyright (c) 2014 Jason Goecke
// messages.go

package wit

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// Message represents a Wit message (https://wit.ai/docs/api#toc_3)
type Message struct {
	MsgID    string    `json:"msg_id"`
	Text     string    `json:"_text"`
	Outcomes []Outcome `json:"outcomes"`
}

// Outcome represents the outcome portion of a Wit message
type Outcome struct {
	Text       string                     `json:"_text"`
	Intent     string                     `json:"intent"`
	IntentId   string                     `json:"intent_id"`
	Entities   map[string][]MessageEntity `json:"entities"`
}

// MessageEntity represents the entity portion of a Wit message
type MessageEntity struct {
	Metadata *string              `json:"metadata,omitempty"`
	Value    *interface{}         `json:"value,omitempty"`
	Grain    *string              `json:"grain,omitempty"`
	Type     *string              `json:"type,omitempty"`
	Unit     *string              `json:"unit,omitempty"`
	Body     *string              `json:"body,omitempty"`
	Entity   *string              `json:"entity,omitempty"`
	Start    *int64               `json:"start,omitempty"`
	End      *int64               `json:"end,omitempty"`
	Values   *[]interface{}       `json:"values,omitempty"`
	From     *DatetimeIntervalEnd `json:"from,omitempty"`
	To       *DatetimeIntervalEnd `json:"to,omitempty"`
	Confidence float32                    `json:"confidence"`
}

// DatetimeValue represents the datetime value portion of a Wit message
type DatetimeValue struct {
	From DatetimeIntervalEnd `json:"from"`
	To   DatetimeIntervalEnd `json:"to"`
}

type DatetimeIntervalEnd struct {
	Value string `json:"value"`
	Grain string `json:"grain"`
}

// MessageRequest represents a request to process a message
type MessageRequest struct {
	File         string `json:"file,omitempty"`
	Query        string `json:"query"`
	MsgID        string `json:"msg_id,omitempty"`
	Context      string `json:"context, omitempty"`
	ContentType  string `json:"contentType, omitempty"`
	N            int    `json:"n,omitempty"`
	FileContents []byte `json:"-"`
	// Are context and Meta necessary anymore?
	// Context     Context
	// Meta        map[string]interface{}
}

// Context represents the context portion of the message request
type Context struct {
	ReferenceTime string `json:"reference_time"`
	Timezone      string `json:"timezone"`
}

// Messages lists an already existing message (https://wit.ai/docs/api#toc_11)
//
//		result, err := client.Messages("ba0fcf60-44d3-4499-877e-c8d65c239730")
func (client *Client) Messages(id string) (*Message, error) {
	result, err := get(client.APIBase + "/messages/" + id)
	if err != nil {
		return nil, err
	}
	message, err := parseMessage(result)
	return message, nil
}

// Message requests processing of a text message (https://wit.ai/docs/api#toc_3)
//
//		result, err := client.Message(request)
func (client *Client) Message(request *MessageRequest) (*Message, error) {
	query := url.QueryEscape(request.Query)
	if request.Context != "" {
		query += "&context=" + request.Context
	}
	if request.MsgID != "" {
		query += "&msg_id" + request.MsgID
	}
	if request.N != 0 {
		query += "&n=" + strconv.Itoa(request.N)
	}
	result, err := get(client.APIBase + "/message?q=" + query)
	if err != nil {
		return nil, err
	}
	message, _ := parseMessage(result)
	return message, nil
}

// AudioMessage requests processing of an audio message (https://wit.ai/docs/api#toc_8)
//
// 		request := &MessageRequest{}
// 		request.File = "./audio_sample/helloWorld.wav"
//		request.FileContents = data
//		request.ContentType = "audio/wav;rate=8000"
// 		message, err := client.AudioMessage(request)
func (client *Client) AudioMessage(request *MessageRequest) (*Message, error) {
	result, err := postFile(client.APIBase+"/speech", request)
	if err != nil {
		return nil, err
	}
	message, err := parseMessage(result)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// Parses the JSON into a Message
//
//		message, err := parseMessage([]byte(data))
func parseMessage(data []byte) (*Message, error) {
	message := &Message{}
	err := json.Unmarshal(data, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}
