package main

import (
	"../."
	"encoding/json"
	"log"
	"os"
)

func main() {
	client := wit.NewClient(os.Getenv("WIT_ACCESS_TOKEN"))

	result, err := client.Intents()
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	log.Println(result)

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	log.Println(string(jsonData[:]))
}

// Output:

// struct:
// &[{52bab837-c31a-46f4-bd74-19a732311ffd name name } {52bab837-2121-4981-a604-61c0809499f5 hello hello } {52bab837-9813-454d-af1d-1c8afd7e4c7c good_bye Good bye } {52bab837-f45c-4707-8763-50d9c337b0ad time time }]

// json:
// [
//     {
//         "id": "52bab837-c31a-46f4-bd74-19a732311ffd",
//         "name": "name",
//         "doc": "name",
//         "metadata": ""
//     },
//     {
//         "id": "52bab837-2121-4981-a604-61c0809499f5",
//         "name": "hello",
//         "doc": "hello",
//         "metadata": ""
//     },
//     {
//         "id": "52bab837-9813-454d-af1d-1c8afd7e4c7c",
//         "name": "good_bye",
//         "doc": "Good bye",
//         "metadata": ""
//     },
//     {
//         "id": "52bab837-f45c-4707-8763-50d9c337b0ad",
//         "name": "time",
//         "doc": "time",
//         "metadata": ""
//     }
// ]
