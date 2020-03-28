package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
)

var (
	output = flag.String("o", "/dev/stdout", "output")
	gofmt  = flag.Bool("f", false, "format")
)

type stringValue struct {
	Value string `xml:"value,attr"`
}

type restriction struct {
	Base         string        `xml:"base,attr"`
	Enumerations []stringValue `xml:"enumeration"`
}

type simpleType struct {
	Name        string      `xml:"name,attr"`
	Restriction restriction `xml:"restriction"`
}

type attribute struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type complexType struct {
	Name       string      `xml:"name,attr"`
	Attributes []attribute `xml:"attribute"`
}

type element struct {
	Name     string `xml:"name,attr"`
	Type     string `xml:"type,attr"`
	Abstract bool   `xml:"abstract,attr"`
}

type xsd struct {
	SimpleTypes  []simpleType  `xml:"simpleType"`
	ComplexTypes []complexType `xml:"complexType"`
	Elements     []element     `xml:"element"`
}

var outputTemplate = template.Must(template.New("output").Funcs(sprig.HermeticTxtFuncMap()).Parse(`package kml

import (
	"image/color"
	"time"
)

{{ range .SimpleTypes -}}
{{ if eq .Restriction.Base "string" -}}
{{ $typeName := .Name | trimSuffix "Type" | title -}}
// A{{ if hasPrefix "A" $typeName }}n{{ end }} {{ $typeName }} is a{{ if hasPrefix "a" .Name }}n{{ end }} {{ .Name }}.
type {{ $typeName }} string

// {{ $typeName }}s.
const (
{{- range .Restriction.Enumerations }}
	{{ $typeName | trimSuffix "Enum" }}{{ .Value | title }} {{ $typeName }} = "{{ .Value }}"
{{- end }}
)
{{ end -}}
{{ end -}}

{{ range .Elements -}}
{{ if and (not .Abstract) (not (regexMatch "^(coordinates|kml|linkSnippet|Scale|Schema|SchemaData|SimpleData|SimpleField|Snippet)$" .Name)) -}}
{{ $functionName := .Name | title | replace "Fov" "FOV" | replace "Http" "HTTP" | replace "Kml" "KML" | replace "Lod" "LOD" | replace "Url" "URL" -}}
{{ $goType := .Type -}}
{{ $constructorName := printf "newSE%s" (.Type | title) -}}
{{ $returnType := "*SimpleElement" -}}
{{ $arg := "value" -}}
{{ $value := "value" -}}
{{ if eq .Name (.Name | title) -}}
	{{ $goType = "...Element" -}}
	{{ $constructorName = "newCE" -}}
	{{ $returnType = "*CompoundElement" -}}
	{{ $arg = "children" -}}
	{{ $value = "children" -}}
{{ else if eq .Type "kml:itemIconStateType" -}}
	{{ $goType = "ItemIconStateEnum" -}}
	{{ $constructorName = "newSEString" -}}
	{{ $value = printf "string(%s)" $arg -}}
{{ else if or (hasSuffix "EnumType" .Type) (eq .Type "kml:itemIconStateType") -}}
	{{ $goType = .Type | trimPrefix "kml:" | trimSuffix "Type" | title -}}
	{{ $constructorName = "newSEString" -}}
	{{ $value = printf "string(%s)" $arg -}}
{{ else if eq .Type "boolean" -}}
	{{ $goType = "bool" -}}
	{{ $constructorName = "newSEBool" -}}
{{ else if or (eq .Type "double") (hasPrefix "kml:angle" .Type) -}}
	{{ $goType = "float64" -}}
	{{ $constructorName = "newSEFloat" -}}
{{ else if eq .Type "anyURI" -}}
	{{ $goType = "string" -}}
	{{ $constructorName = "newSEString" -}}
{{ else if eq .Type "kml:colorType" -}}
	{{ $goType = "color.Color" -}}
	{{ $constructorName = "newSEColor" -}}
{{ else if eq .Type "kml:dateTimeType" -}}
	{{ $goType = "time.Time" -}}
	{{ $constructorName = "newSETime" -}}
{{ else if eq .Type "kml:vec2Type" -}}
	{{ $goType = "Vec2" -}}
	{{ $constructorName = "newSEVec2" -}}
{{ else if or (eq .Type "kml:BoundaryType") (eq .Type "kml:KmlType") (eq .Type "kml:SnippetType") -}}
	{{ $goType = "Element" -}}
	{{ $constructorName = "newSEElement" -}}
	{{ $returnType = "*CompoundElement" -}}
{{ end -}}
// {{ $functionName }} returns a new {{ .Name }} element.
func {{ $functionName }}({{ $arg }} {{ $goType }}) {{ $returnType }} {
	return {{ $constructorName }}("{{ .Name }}", {{ $value }})
}
{{ end -}}
{{ end -}}
`))

func run() error {
	flag.Parse()

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		return err
	}
	defer f.Close()

	xsd := &xsd{}
	if err := xml.NewDecoder(f).Decode(xsd); err != nil {
		return err
	}

	source := &bytes.Buffer{}
	if err := outputTemplate.Execute(source, xsd); err != nil {
		return err
	}

	if !*gofmt {
		return ioutil.WriteFile(*output, source.Bytes(), 0666)
	}

	formattedSource, err := format.Source(source.Bytes())
	if err != nil {
		return err
	}

	return ioutil.WriteFile(*output, formattedSource, 0666)
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
