// Copyright (c) 2014 Jason Goecke
// messages.go

package wit

import (
	"encoding/json"
	"net/url"
)

// Represents a Wit message (https://wit.ai/docs/api#toc_3)
type Message struct {
	MsgId   string  `json:"msg_id"`
	MsgBody string  `json:"msg_body"`
	Outcome Outcome `json:"outcome"`
}

// Represents the outcome portion of a Wit message
type Outcome struct {
	Intent     string        `json:"intent"`
	Entities   MessageEntity `json:"entities"`
	Confidence float32       `json:"confidence"`
}

// Represents the entity portion of a Wit message
type MessageEntity struct {
	Metric   Metric     `json:"metric"`
	Datetime []Datetime `json:"datetime"`
}

// Represents the metric portion of a Wit message
type Metric struct {
	Value string `json:"value"`
	Body  string `json:"value"`
}

// Represents the datetime portion of a Wit message
type Datetime struct {
	Value DatetimeValue `json:"value"`
	Body  string        `json:"body"`
}

// Represents the datetime value portion of a Wit message
type DatetimeValue struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Represents a request to process a message
type MessageRequest struct {
	File, Query, MsgId, ContentType string
	// Are context and Meta necessary anymore?
	// Context     Context
	// Meta        map[string]interface{}
}

// Represents the context portion of the message request
type Context struct {
	ReferenceTime string `json:"reference_time"`
	Timezone      string `json:"timezone"`
}

// Lists an already existing message (https://wit.ai/docs/api#toc_11)
//
//		result, err := client.Messages("ba0fcf60-44d3-4499-877e-c8d65c239730")
func (client *WitClient) Messages(id string) (*Message, error) {
	result, err := get(client.ApiBase + "/messages/" + id)
	if err != nil {
		return nil, err
	}
	message, err := parseMessage(result)
	return message, nil
}

// Requests processing of a text message (https://wit.ai/docs/api#toc_3)
//
//		result, err := client.Message(request)
func (client *WitClient) Message(request *MessageRequest) (*Message, error) {
	result, err := get(client.ApiBase + "/message?q=" + url.QueryEscape(request.Query))
	if err != nil {
		return nil, err
	}
	message, _ := parseMessage(result)
	return message, nil
}

// Requests processing of an audio message (https://wit.ai/docs/api#toc_8)
//
// 		request := &MessageRequest{}
// 		request.File = "./audio_sample/helloWorld.wav"
//		request.ContentType = "audio/wav;rate=8000"
// 		message, err := client.AudioMessage(request)
func (client *WitClient) AudioMessage(request *MessageRequest) (*Message, error) {
	result, err := postFile(client.ApiBase+"/speech", request)
	if err != nil {
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
