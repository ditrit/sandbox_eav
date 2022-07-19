package utils

import "strings"

// build a json from a string slices
// turn []string {`"a",1`, `"b",1`} to `{"a":1,"b":2}`
func BuildJsonFromStrings(pairs []string) string {
	var b strings.Builder
	// starting the json string
	b.WriteString("{")
	b.WriteString(strings.Join(pairs, ","))
	//ending the json string
	b.WriteString("}")
	return b.String()
}
