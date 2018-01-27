// hike-and-fly-route prints a KML file of the route of popular races.
package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/twpayne/go-kml"
	"github.com/twpayne/go-kml/icon"
	"github.com/twpayne/go-kml/sphere"
)

var (
	raceFlag = flag.String("race", "x-pyr-2016", "race")
)

type turnpoint struct {
	name   string
	lat    float64
	lon    float64
	radius int
	paddle string
	notes  string
}

type race struct {
	name       string
	snippet    string
	turnpoints []turnpoint
}

var (
	races = map[string]race{
		"x-pyr-2016": {
			name:    "X-Pyr 2016",
			snippet: "http://www.x-pyr.com/",
			turnpoints: []turnpoint{
				{
					name:   "Start: Hondarribia",
					lat:    43.378709,
					lon:    -1.796020,
					paddle: "go",
					notes:  "Beginning of timed section.",
				},
				{
					name:   "TP1: Larun",
					lat:    43.309111,
					lon:    -1.635409,
					paddle: "1",
					notes:  "Pilots must walk across the gate on the summit. If it is flyable, the organisation will mark the take off.",
				},
				{
					name:   "TP2: Orhi",
					lat:    42.988113,
					lon:    -1.005943,
					radius: 100,
					paddle: "2",
				},
				{
					name:   "TP3: Anayet",
					lat:    42.781424,
					lon:    -0.455415,
					radius: 400,
					paddle: "3",
				},
				{
					name:   "TP4: Peña Montañesa",
					lat:    42.490226,
					lon:    0.199227,
					radius: 1000,
					paddle: "4",
				},
				{
					name:   "TP5: Ceciré",
					lat:    42.757425,
					lon:    0.537662,
					radius: 400,
					paddle: "5",
				},
				{
					name:   "TP6: Berguedà",
					lat:    42.248419,
					lon:    1.885515,
					paddle: "6",
					notes:  "Pilots must walk across the gate on the waypoint.",
				},
				{
					name:   "TP7: Canigó",
					lat:    42.519159,
					lon:    2.456149,
					radius: 3000,
					paddle: "7",
				},
				{
					name:   "TP8: Santa Helena de Rodes",
					lat:    42.326468,
					lon:    3.16018,
					paddle: "stop",
					notes:  "End of timed section. Pilots must cross a signposted area on foot.",
				},
				{
					name:   "Finish: El Port de la Selva",
					lat:    42.336152,
					lon:    3.201039,
					paddle: "ylw-stars",
				},
			},
		},
	}
)

func (tp turnpoint) folder() kml.Element {
	center := kml.Coordinate{Lon: tp.lon, Lat: tp.lat}
	var radiusPlacemark kml.Element
	if tp.radius != 0 {
		radiusPlacemark = kml.Placemark(
			kml.LineString(
				kml.Coordinates(sphere.FAI.Circle(center, float64(tp.radius), 1)...),
			),
			kml.Style(
				kml.LineStyle(
					kml.Color(color.RGBA{R: 0, G: 192, B: 0, A: 192}),
				),
			),
		)
	}
	var snippet kml.Element
	switch {
	case tp.notes != "":
		snippet = kml.Snippet(tp.notes)
	case tp.radius != 0:
		snippet = kml.Snippet(fmt.Sprintf("%dm radius.", tp.radius))
	}
	return kml.Folder(
		kml.Name(tp.name),
		snippet,
		kml.Placemark(
			kml.Point(
				kml.Coordinates(center),
			),
			kml.Style(
				icon.PaddleIconStyle(tp.paddle),
			),
		),
		radiusPlacemark,
		kml.Style(
			kml.ListStyle(
				kml.ListItemType("checkHideChildren"),
			),
		),
	)
}

func (r race) routeFolder() kml.Element {
	var coordinates []kml.Coordinate
	for _, tp := range r.turnpoints {
		coordinates = append(coordinates, kml.Coordinate{Lon: tp.lon, Lat: tp.lat})
	}
	return kml.Folder(
		kml.Name("Route"),
		kml.Placemark(
			kml.LineString(
				kml.Coordinates(coordinates...),
				kml.Tessellate(true),
			),
			kml.Style(
				kml.LineStyle(
					kml.Color(color.RGBA{R: 144, G: 144, B: 0, A: 192}),
					kml.Width(4),
				),
			),
		),
		kml.Style(
			kml.ListStyle(
				kml.ListItemType("checkHideChildren"),
			),
		),
	)
}

func (r race) document() kml.Element {
	var folders []kml.Element
	folders = append(folders, r.routeFolder())
	for _, tp := range r.turnpoints {
		folders = append(folders, tp.folder())
	}
	return kml.KML(
		kml.Document(append([]kml.Element{
			kml.Name(fmt.Sprintf("%s Route", r.name)),
			kml.Snippet(r.snippet),
			kml.Open(true),
		}, folders...)...),
	)
}

func run() error {
	r, ok := races[*raceFlag]
	if !ok {
		return fmt.Errorf("unknown race: %q", *raceFlag)
	}
	return r.document().WriteIndent(os.Stdout, "", "  ")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
