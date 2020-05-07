package mongoke

import "encoding/json"

func prettyPrint(x ...interface{}) {
	for _, x := range x {
		json, err := json.MarshalIndent(x, "", "   ")
		if err != nil {
			panic(err)
		}
		println(string(json))
	}
	println()
}
