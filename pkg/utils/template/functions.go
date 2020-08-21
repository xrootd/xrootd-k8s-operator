package template

import "text/template"

// TemplateFunctions provides helper functions to be used
// in go templates.
var TemplateFunctions = template.FuncMap{
	"Iterate": iterateCount,
}

// iterateCount returns array filled with 0..count
func iterateCount(count int) []int {
	items := make([]int, count)
	for i := 0; i < count; i++ {
		items[i] = i
	}
	return items
}
