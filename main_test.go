package main

import (
	"testing"
)

func TestAll(t *testing.T) {
	var template string = `
Packages: {{len .Packages}}
{{range $k, $v := .Packages}}
package {{$k}}
  {{(index (index (index $v.Files "main.go").Comments 0).List 0).Text}}
{{end}}
`
	var test string = `
Packages: 1

package main
  //TEST COMMENT - DO NOT REMOVE

`
	result := render(template, parseDir("."))
	t.Log(result)
	if result != test {
		t.Fail()
	}
}
