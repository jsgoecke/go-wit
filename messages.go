// Copyright (c) 2014 Jason Goecke
// messages.go

package wit

import (
	"encoding/json"
	"net/url"
)

// Message represents a Wit message (https://wit.ai/docs/api#toc_3)
type Message struct {
	MsgID   string  `json:"msg_id"`
	MsgBody string  `json:"msg_body"`
	Outcome Outcome `json:"outcome"`
}

// Outcome represents the outcome portion of a Wit message
type Outcome struct {
	Intent     string        `json:"intent"`
	Entities   MessageEntity `json:"entities"`
	Confidence float32       `json:"confidence"`
}

// MessageEntity represents the entity portion of a Wit message
type MessageEntity struct {
	Metric   Metric     `json:"metric"`
	Datetime []Datetime `json:"datetime"`
}

// Metric represents the metric portion of a Wit message
type Metric struct {
	Value string `json:"value"`
	Body  string `json:"value"`
}

// Datetime represents the datetime portion of a Wit message
type Datetime struct {
	Value DatetimeValue `json:"value"`
	Body  string        `json:"body"`
}

// DatetimeValue represents the datetime value portion of a Wit message
type DatetimeValue struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// MessageRequest represents a request to process a message
type MessageRequest struct {
	File, Query, MsgID, ContentType string
	FileContents                    []byte
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
	result, err := get(client.APIBase + "/message?q=" + url.QueryEscape(request.Query))
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
