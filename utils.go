package mongoke

import "encoding/json"

func prettyPrint(x interface{}) {
	json, err := json.MarshalIndent(x, "", "   ")
	if err != nil {
		panic(err)
	}
	println(string(json))
}
