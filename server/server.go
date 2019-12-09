package main

import (
	"html/template"
	"io"
	"log"

	"github.com/donniet/tictacgo"
)

const (
	positionHTMLTemplateString = `<table rows="3" cols="3"><tbody>
		<tr><td>{{ .Get(0,0) }}</td><td>{{ .Get(1,0) }}</td><td>{{ .Get(2,0) }}</td></tr>
		<tr><td>{{ .Get(0,1) }}</td><td>{{ .Get(1,1) }}</td><td>{{ .Get(2,1) }}</td></tr>
		<tr><td>{{ .Get(0,2) }}</td><td>{{ .Get(1,2) }}</td><td>{{ .Get(2,2) }}</td></tr>
	</tbody></table>`

	pageHTMLTemplateString = `
<!doctype html>
<html>
<head><title>Tic Tac</title></head>
<body>

</body>
</html>`
)

var (
	positionHTMLTemplate *template.Template
)

func init() {
	var err error
	if positionHTMLTemplate, err = template.New("position").Parse(positionHTMLTemplateString); err != nil {
		log.Fatalf("position template does not compile: %s", err)
	}

}

// WritePosition writes out the position using the embedded template to a Writer
func WritePosition(w io.Writer, p tictacgo.Position) {
	positionHTMLTemplate.ExecuteTemplate(w, "position", p)
}
