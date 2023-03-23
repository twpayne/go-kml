package main

import (
	_ "embed"
	"encoding/xml"
	"flag"
	"fmt"
	"go/format"
	"os"
	"regexp"
	"strings"
	"text/template"
	"unicode"
)

var (
	output    = flag.String("o", "/dev/stdout", "output")
	namespace = flag.String("n", "", "namespace")

	aOrAnRegexp         = regexp.MustCompile(`(?i)\A[aeio]`)
	abbreviationsRegexp = regexp.MustCompile(`Fov|Http|Id|Kml|Lod|Url`)
	gxEnumTypeRegexp    = regexp.MustCompile(`\Agx:(.*Enum)Type\z`)
	kmlEnumTypeRegexp   = regexp.MustCompile(`\Akml:(.*Enum)Type\z`)

	xsdTypeToGoType = map[string]string{
		"anyURI":                "string",
		"boolean":               "bool",
		"double":                "float64",
		"float":                 "float64",
		"gx:outerWidthType":     "float64",
		"int":                   "int",
		"integer":               "int",
		"kml:angle180Type":      "float64",
		"kml:angle360Type":      "float64",
		"kml:angle90Type":       "float64",
		"kml:anglepos180Type":   "float64",
		"kml:colorType":         "color.Color",
		"kml:dateTimeType":      "time.Time",
		"kml:itemIconStateType": "ItemIconStateEnum",
		"kml:SchemaDataType":    "string",
		"kml:SimpleDataType":    "string",
		"kml:vec2Type":          "Vec2",
		"string":                "string",
	}
)

//go:embed output.go.tmpl
var outputGoTemplateText string

func run() error {
	flag.Parse()

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		return err
	}
	defer file.Close()

	var schema Schema
	if err := xml.NewDecoder(file).Decode(&schema); err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"aOrAn": func(s string) string {
			if aOrAnRegexp.MatchString(s) {
				return "an " + s
			}
			return "a " + s
		},
		"hasSuffix": func(suffix, s string) bool {
			return strings.HasSuffix(s, suffix)
		},
		"nameToGoName": func(name string) string {
			return abbreviationsRegexp.ReplaceAllStringFunc(name, strings.ToUpper)
		},
		"titleFirst": titleFirst,
		"trimSuffix": func(suffix, s string) string {
			return strings.TrimSuffix(s, suffix)
		},
		"xsdTypeToGoType": func(typeStr string) string {
			if goType, ok := xsdTypeToGoType[typeStr]; ok {
				return goType
			}
			if match := kmlEnumTypeRegexp.FindStringSubmatch(typeStr); match != nil {
				return titleFirst(match[1])
			}
			if match := gxEnumTypeRegexp.FindStringSubmatch(typeStr); match != nil {
				return "Gx" + titleFirst(match[1])
			}
			return ""
		},
	}

	outputGoTemplate, err := template.New("output.go.tmpl").
		Funcs(funcMap).
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

	formattedSource, err := format.Source([]byte(source.String()))
	if err != nil {
		formattedSource = []byte(source.String())
	}

	return os.WriteFile(*output, formattedSource, 0o666)
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func titleFirst(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToTitle(runes[0])
	return string(runes)
}
