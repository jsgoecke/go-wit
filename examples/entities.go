package main

import (
	"../."
	"encoding/json"
	"log"
	"os"
)

func main() {
	client := wit.NewClient(os.Getenv("WIT_ACCESS_TOKEN"))

	// Show a Wit builtin entity
	entity, err := client.Entity("wit$age_of_person")
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	log.Println(entity)
	jsonData, _ := json.MarshalIndent(entity, "", "    ")
	log.Println(string(jsonData[:]))

	// Show a custom entity
	entity, err = client.Entity("68e5fbfb-1839-422a-8e52-28798344b2ad$favorite_city")
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	log.Println(entity)
	jsonData, _ = json.MarshalIndent(entity, "", "    ")
	log.Println(string(jsonData[:]))
}

// Output:

//structs:
// &{true (BETA) The age (in years) of a person, like in `22 years old`. The value of the entity is an integer. Does not support smaller granularity (months, weeks, etc.). In the expression, just select the integer value (`22`). wit$age_of_person []}
// &{false These are cities worth going to 68e5fbfb-1839-422a-8e52-28798344b2ad$favorite_city [{Barcelona [Gaudi Med Paella Sagrada Familia]}]}

//json:
// {
//     "builtin": true,
//     "doc": "(BETA) The age (in years) of a person, like in `22 years old`. The value of the entity is an integer. Does not support smaller granularity (months, weeks, etc.). In the expression, just select the integer value (`22`).",
//     "id": "wit$age_of_person",
//     "Values": null
// }
// {
//     "builtin": false,
//     "doc": "These are cities worth going to",
//     "id": "68e5fbfb-1839-422a-8e52-28798344b2ad$favorite_city",
//     "Values": [
//         {
//             "value": "Barcelona",
//             "expressions": [
//                 "Gaudi",
//                 "Med",
//                 "Paella",
//                 "Sagrada Familia"
//             ]
//         }
//     ]
// }
