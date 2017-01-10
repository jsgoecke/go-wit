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
	  "msg_id" : "2f41839e-2b54-4de2-aa59-fc016c3e58d1",
	  "_text" : "how many people between Tuesday and Friday?",
	  "outcomes" : [ {
	    "_text" : "how many people between Tuesday and Friday?",
	    "intent" : "query_metrics",
	    "entities" : {
        "intent": [
	        "datetime" : [ {
	          "type" : "interval",
	          "from" : {
	            "value" : "2015-12-01T00:00:00.000-08:00",
	            "grain" : "day"
	          },
	          "to" : {
	            "value" : "2015-12-05T00:00:00.000-08:00",
	            "grain" : "day"
	          },
	          "values" : [ {
	            "type" : "interval",
	            "from" : {
	              "value" : "2015-12-01T00:00:00.000-08:00",
	              "grain" : "day"
	            },
	            "to" : {
	              "value" : "2015-12-05T00:00:00.000-08:00",
	              "grain" : "day"
	            }
	          }, {
	            "type" : "interval",
	            "from" : {
	              "value" : "2015-12-08T00:00:00.000-08:00",
	              "grain" : "day"
	            },
	            "to" : {
	              "value" : "2015-12-12T00:00:00.000-08:00",
	              "grain" : "day"
	            }
	          }, {
	            "type" : "interval",
	            "from" : {
	              "value" : "2015-12-15T00:00:00.000-08:00",
	              "grain" : "day"
	            },
	            "to" : {
	              "value" : "2015-12-19T00:00:00.000-08:00",
	              "grain" : "day"
	            }
	          } ]
	        } ],
	      "confidence" : 0.522
      ] },
	  } ]
	}`

	message, err := parseMessage([]byte(data))
	if err != nil {
		t.Error(err.Error())
	}

	if message.MsgID != "2f41839e-2b54-4de2-aa59-fc016c3e58d1" {
		t.Errorf("not equal %s != %s", "2f41839e-2b54-4de2-aa59-fc016c3e58d1", message.MsgID)
	}
	if message.Outcomes[0].Intent != "query_metrics" {
		t.Errorf("not equal %s != %s", "query_metrics", message.Outcomes[0].Intent)
	}
	if message.Outcomes[0].Entities["datetime"][0].From.Grain != "day" {
		t.Errorf("not equal %s != %s", "day", message.Outcomes[0].Entities["datetime"][0].From.Grain)
	}
}

func TestWitMessageRequest(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))

	request := &MessageRequest{}
	request.Query = "Hello world"
	result, err := client.Message(request)
	if err != nil {
		t.Error(err)
		return
	}
	if result.Text != "Hello world" {
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
		if message.Text != "hello world" {
			t.Error("Audio POST did not work properly")
		}
	}
}

func TestWitPostAudioContentsMessage(t *testing.T) {
	client := NewClient(os.Getenv("WIT_ACCESS_TOKEN"))
	file, err := os.Open("./audio_sample/helloWorld.wav")
	if err != nil {
		t.Error(err)
		return
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
		if message.Text != "hello world" {
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
		t.Error(err)
	} else {
		if message.MsgID != msgID {
			t.Error("Message JSON did not parse properly.")
		}
	}
}
