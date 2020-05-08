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
func pretty(x ...interface{}) string {
	res := ""
	for _, x := range x {
		json, err := json.MarshalIndent(x, "", "   ")
		if err != nil {
			panic(err)
		}
		res += string(json)
		res += "\n"
	}
	return res
}
