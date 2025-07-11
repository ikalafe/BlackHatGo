package main

import (
	"os"
	"text/template"
)

var x = `
<html>
	<body>
		Hello{{.}}
	</body>
</html>
`

func main() {
	t, err := template.New("hello").Parse(x)
	if err != nil {
		panic(err)
	}

	t.Execute(os.Stdout, "<script>alert('world')</script>")
}