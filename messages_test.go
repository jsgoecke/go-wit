// Copyright (c) 2014 Jason Goecke
// messages_test.go

package wit

import (
	"os"
	"testing"
	"time"
)

var msgID string

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

	if message.MsgID != "1234" ||
		message.Outcome.Intent != "query_metrics" ||
		message.Outcome.Entities.Datetime[0].Body != "Tuesday" {
		t.Error("Message JSON did not parse properly.")
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
	msgID = result.MsgID
}

func TestWitPostAudioMessage(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	request := &MessageRequest{}
	request.File = "./audio_sample/helloWorld.wav"
	request.ContentType = "audio/wav"
	message, err := client.AudioMessage(request)
	if err != nil {
		t.Error(err)
	} else {
		if message.MsgBody != "hello world" {
			t.Error("Audio POST did not work properly")
		}
	}
}

func TestWitPostAudioContentsMessage(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	file, err := os.Open("./audio_sample/helloWorld.wav")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	stats, statsErr := file.Stat()
	if statsErr != nil {
		t.Error(err)
	}
	size := stats.Size()
	data := make([]byte, size)
	file.Read(data)
	request := &MessageRequest{}
	request.FileContents = data
	request.ContentType = "audio/wav"
	message, err := client.AudioMessage(request)
	if err != nil {
		t.Error(err)
	} else {
		if message.MsgBody != "hello world" {
			t.Error("Audio POST did not work properly")
		}
	}
}

func TestWitMessages(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	//Wait for the message to be indexed
	time.Sleep(300 * time.Millisecond)
	message, err := client.Messages(msgID)
	if err != nil {
		t.Error("Message JSON did not parse properly.")
	} else {
		if message.MsgID != msgID {
			t.Error("Message JSON did not parse properly.")
		}
	}
}
