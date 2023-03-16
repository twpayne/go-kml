package main

import (
	_ "embed"
	"encoding/xml"
	"flag"
	"fmt"
	"go/format"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

var (
	output    = flag.String("o", "/dev/stdout", "output")
	gofmt     = flag.Bool("f", false, "format")
	namespace = flag.String("n", "", "namespace")
)

//go:embed output.go.tmpl
var outputGoTemplateText string

func run() error {
	flag.Parse()

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		return err
	}
	defer f.Close()

	var schema Schema
	if err := xml.NewDecoder(f).Decode(&schema); err != nil {
		return err
	}

	outputGoTemplate, err := template.New("").
		Funcs(sprig.HermeticTxtFuncMap()).
		Parse(outputGoTemplateText)
	if err != nil {
		return err
	}

	source := &strings.Builder{}
	if err := outputGoTemplate.Execute(source, struct {
		Namespace string
		Schema    Schema
	}{
		Namespace: *namespace,
		Schema:    schema,
	}); err != nil {
		return err
	}

	if !*gofmt {
		return os.WriteFile(*output, []byte(source.String()), 0o666)
	}

	formattedSource, err := format.Source([]byte(source.String()))
	if err != nil {
		return err
	}

	return os.WriteFile(*output, formattedSource, 0o666)
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
