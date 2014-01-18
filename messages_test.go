// Copyright (c) 2014 Jason Goecke
// messages_test.go

package wit

import (
	"os"
	"testing"
)

func TestWitMessageParsing(t *testing.T) {
	data := `
	{
	   "msg_id":"1234",
	   "msg_body":"how many people between Tuesday and Friday?",
	   "outcome":{
	      "intent":"query_metrics",
	      "entities":{
	         "metric":{
	            "value":"metric_visitors",
	            "body":"people"
	         },
	         "datetime":[
	            {
	               "value":{
	                  "from":"2013-10-21T00:00:00.000Z",
	                  "to":"2013-10-22T00:00:00.000Z"
	               },
	               "body":"Tuesday"
	            },
	            {
	               "value":{
	                  "from":"2013-10-24T00:00:00.000Z",
	                  "to":"2013-10-25T00:00:00.000Z"
	               },
	               "body":"Friday"
	            }
	         ]
	      },
	      "confidence":0.979
	   }
	}`

	message, err := parseMessage([]byte(data))
	if err != nil {
		t.Error(err.Error())
	}

	if message.MsgId != "1234" ||
		message.Outcome.Intent != "query_metrics" ||
		message.Outcome.Entities.Datetime[0].Body != "Tuesday" {
		t.Error("Message JSON did not parse properly.")
	}
}

func TestWitMessages(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	message, err := client.Messages("7c78b1c7-4845-4d69-8bb7-01854fa2b792")
	if err != nil {
		t.Error("Message JSON did not parse properly.")
	} else {
		if message.MsgId != "7c78b1c7-4845-4d69-8bb7-01854fa2b792" {
			t.Error("Message JSON did not parse properly.")
		}
	}
}

func TestWitMessageRequest(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))

	request := &MessageRequest{}
	request.Query = "Hello world"
	result, err := client.Message(request)
	if err != nil {
		t.Error(err)
	}
	if result.MsgBody != "Hello world" {
		t.Error("Did not process properly")
	}
}

func TestWitPostAudioMessage(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	request := &MessageRequest{}
	request.File = "./audio_sample/helloWorld.wav"
	request.ContentType = "audio/wav;rate=8000"
	message, err := client.AudioMessage(request)
	if err != nil {
		t.Error(err)
	} else {
		if message.MsgBody != "hello world" {
			t.Error("Audio POST did not work properly")
		}
	}
}
