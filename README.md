#go-wit 

[![wercker status](https://app.wercker.com/status/db550f8016d3e02bd1d79d94bccf764b "wercker status")](https://app.wercker.com/project/bykey/db550f8016d3e02bd1d79d94bccf764b)

A Go library for the [Wit.ai](http://wit.ai) API for Natural Language Processing.

## Version

0.3 

## Installation

	go get github.com/jsgoecke/go-wit

## Documentation

[http://godoc.org/github.com/jsgoecke/go-wit](http://godoc.org/github.com/jsgoecke/go-wit)

## Usage

```go
package main

import (
	"github.com/jsgoecke/go-wit"
	"encoding/json"
	"log"
	"os"
)

func main() {
	client := wit.NewClient(os.Getenv("WIT_ACCESS_TOKEN"))

	// Process a text message
	request := &wit.MessageRequest{}
	request.Query = "Hello world"
	result, err := client.Message(request)
	if err != nil {
		println(err)
		os.Exit(-1)
	}
	log.Println(result)
	data, _ := json.MarshalIndent(result, "", "    ")
	log.Println(string(data[:]))

	// Process an audio/wav message
	request = &wit.MessageRequest{}
	request.File = "../audio_sample/helloWorld.wav"
	request.ContentType = "audio/wav;rate=8000"
	result, err = client.AudioMessage(request)
	if err != nil {
		println(err)
		os.Exit(-1)
	}
	log.Println(result)
	data, _ = json.MarshalIndent(result, "", "    ")
	log.Println(string(data[:]))
}

// Output:

// structs:
// &{bf699a8f-bc90-4fb4-a715-bd8bd77749db Hello world {hello {{ } []} 0.996}}
// &{54ed4e6d-0653-453e-8c0c-81da57c3846c hello world {hello {{ } []} 0.993}}

// json:
// {
//     "msg_id": "76f1c370-bd92-417f-8cb3-e1419d1a9cb3",
//     "_text": "Hello world",
//     "outcome": {
//         "_text": "Hello world",
//         "intent": "default_intent",
//         "intent_id": "",
//         "entities": {
//             "intent": [
//                 {
//                     "value": "greetings",
//                     "confidence": 0.5108957
//                 }
//             ]
//         }
//      }
// }
// {
//     "msg_id": "322f9b61-0f75-4953-a392-f8eca058a12f",
//     "_text": "hello world",
//     "outcome": {
//         "intent": "default_intent",
//         "_text": "hello world",
//         "entities": {
//             "intent": [
//               "value": "greetings",
//               "confidence":  0.5108957
//             ]
//         }
//     }
// }
```

## Testing

Must have the environment variable WIT_ACCESS_TOKEN set to your Wit API token.
	
	cd go-wit
	go test

### Test Coverage

[http://gocover.io/github.com/jsgoecke/go-wit](http://gocover.io/github.com/jsgoecke/go-wit)

### Lint

[http://go-lint.appspot.com/github.com/jsgoecke/go-wit](http://go-lint.appspot.com/github.com/jsgoecke/go-wit)

## License

MIT, see LICENSE.txt

## Author

Jason Goecke [@jsgoecke](http://twitter.com/jsgoecke).
