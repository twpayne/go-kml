// hike-and-fly-route prints a KML file of the route of popular races.
package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"net/url"
	"os"
	"strconv"

	kml "github.com/twpayne/go-kml"
	"github.com/twpayne/go-kml/icon"
	"github.com/twpayne/go-kml/sphere"
	polyline "github.com/twpayne/go-polyline"
)

var (
	formatFlag = flag.String("format", "kml", "format")
	raceFlag   = flag.String("race", "eigertour-2019-challenge", "race")
)

type turnpoint struct {
	name      string
	lat       float64
	lon       float64
	radius    int
	paddle    string
	signboard bool
	pass      string
	offRoute  bool
	notes     string
}

type race struct {
	name       string
	snippet    string
	turnpoints []turnpoint
}

var (
	berghausBaregg   = kml.Coordinate{Lat: 46.60046, Lon: 8.060011}
	berghausBareggNW = sphere.FAI.Offset(berghausBaregg, 100, 315)
	berghausBareggSE = sphere.FAI.Offset(berghausBaregg, 100, 135)

	lobhornhutte   = kml.Coordinate{Lat: 46.618514, Lon: 7.868981}
	lobhornhutteNW = sphere.FAI.Offset(lobhornhutte, 100, 315)
	lobhornhutteSE = sphere.FAI.Offset(lobhornhutte, 100, 135)
)

var races = map[string]race{
	"eigertour-2019-challenge": {
		name:    "Eigertour 2019 Challenge",
		snippet: "https://eigertour.rocks/ Created by twpayne@gmail.com",
		turnpoints: []turnpoint{
			{
				name:   "Eigerplatz",
				lat:    dms(46, 37, 25.4),
				lon:    dms(8, 2, 6.7),
				paddle: "ylw-stars",
			},
			{
				name:   "First",
				lat:    dms(46, 39, 31.5),
				lon:    dms(8, 3, 15.4),
				paddle: "A",
			},
			{
				name:   "Tierberglihütte",
				lat:    46.702018,
				lon:    8.41421,
				paddle: "B",
			},
			{
				name:   "Gaulihütte",
				lat:    46.623778,
				lon:    8.216637,
				paddle: "C",
			},
			{
				name:   "Lobhornhütte",
				lat:    lobhornhutteNW.Lat,
				lon:    lobhornhutteNW.Lon,
				paddle: "D",
			},
			{
				name:   "Niesen",
				lat:    46.644999,
				lon:    7.651387,
				paddle: "E",
			},
			{
				name:   "Doldehornhütte",
				lat:    46.486806,
				lon:    7.697366,
				paddle: "F",
			},
			{
				name:   "Schmadrihütte",
				lat:    46.499159,
				lon:    7.892225,
				paddle: "G",
			},
			{
				name:   "Berghaus Bäregg",
				lat:    berghausBareggSE.Lat,
				lon:    berghausBareggSE.Lon,
				paddle: "H",
			},
			{
				name:   "Glecksteinhütte",
				lat:    46.625129,
				lon:    8.096503,
				paddle: "I",
			},
			{
				name:   "Lobhornhütte",
				lat:    lobhornhutteSE.Lat,
				lon:    lobhornhutteSE.Lon,
				paddle: "J",
			},
			{
				name:   "Berghaus Bäregg",
				lat:    berghausBareggNW.Lat,
				lon:    berghausBareggNW.Lon,
				paddle: "K",
			},
			{
				name: "Eigerplatz",
				lat:  dms(46, 37, 25.4),
				lon:  dms(8, 2, 6.7),
			},
		},
	},
	"red-bull-x-alps-2019": {
		name:    "Red Bull X-Alps 2019",
		snippet: "https://www.redbullxalps.com/ Created by twpayne@gmail.com",
		turnpoints: []turnpoint{
			{
				name:   "Salzburg",
				lat:    47.79885,
				lon:    13.0484,
				paddle: "go",
			},
			{
				name:      "Gaisberg",
				lat:       47.804133,
				lon:       13.110917,
				paddle:    "1",
				signboard: true,
			},
			{
				name:      "Wagrain-Kleinarl",
				lat:       47.332295,
				lon:       13.305787,
				paddle:    "2",
				signboard: true,
			},
			{
				name:      "Aschau-Chiemsee",
				lat:       47.784362,
				lon:       12.33277,
				paddle:    "3",
				signboard: true,
			},
			{
				name:      "Kronplatz",
				lat:       46.737598,
				lon:       11.9549,
				paddle:    "4",
				signboard: true,
			},
			{
				name: "Zugspitz",
				lat:  47.4211,
				lon:  10.98526,
				pass: "N",
			},
			{
				name:      "Lermoos-Tiroler Zugspitz Arena",
				lat:       47.401283,
				lon:       10.879767,
				paddle:    "5",
				signboard: true,
			},
			{
				name:      "Davos",
				lat:       46.815225,
				lon:       9.851879,
				paddle:    "6",
				signboard: true,
			},
			{
				name:      "Titlis",
				lat:       46.770918,
				lon:       8.424457,
				paddle:    "7",
				signboard: true,
			},
			{
				name:   "Eiger",
				lat:    46.577621,
				lon:    8.005393,
				paddle: "8",
				radius: 1500,
			},
			{
				name:     "Mont Blanc",
				lat:      45.830359,
				lon:      6.867674,
				paddle:   "9",
				pass:     "N",
				offRoute: true,
			},
			{
				name:      "St. Hilare",
				lat:       45.306816,
				lon:       5.887857,
				paddle:    "10",
				signboard: true,
			},
			{
				name:   "Monte Viso",
				lat:    44.667312,
				lon:    7.090381,
				paddle: "A",
				radius: 2250,
			},
			{
				name:   "Cheval Blanc",
				lat:    44.120985,
				lon:    6.422229,
				paddle: "B",
				pass:   "W",
			},
			{
				name:      "Peille",
				lat:       43.755956,
				lon:       7.410751,
				paddle:    "stop",
				signboard: true,
			},
			{
				name:   "Monaco",
				lat:    43.75875,
				lon:    7.454787,
				paddle: "ylw-stars",
			},
		},
	},
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
	"x-pyr-2018": {
		name:    "X-Pyr 2018",
		snippet: "http://www.x-pyr.com/",
		turnpoints: []turnpoint{
			{
				name:   "Hondarribia",
				lat:    43.379469,
				lon:    -1.796731,
				paddle: "go",
			},
			{
				name:   "La Rhune",
				lat:    43.309039,
				lon:    -1.635419,
				paddle: "1",
			},
			{
				name:   "Orhi",
				lat:    42.988111,
				lon:    -1.005939,
				paddle: "2",
			},
			{
				name:   "Midi d-Ossau",
				lat:    42.843250,
				lon:    -0.438069,
				paddle: "3",
			},
			{
				name:   "Turbon",
				lat:    42.416931,
				lon:    0.505181,
				paddle: "4",
			},
			{
				name:   "Midi de Bigorre",
				lat:    42.937019,
				lon:    0.140761,
				paddle: "5",
			},
			{
				name:   "Pedraforca",
				lat:    42.239869,
				lon:    1.702950,
				paddle: "6",
			},
			{
				name:   "Canigo",
				lat:    42.519161,
				lon:    2.456150,
				paddle: "7",
			},
			{
				name:   "Santa Helena de Rodes",
				lat:    42.326469,
				lon:    3.160181,
				paddle: "stop",
			},
			{
				name:   "El Port de la Selva",
				lat:    42.336150,
				lon:    3.201039,
				paddle: "ylw-stars",
			},
		},
	},
}

func (tp turnpoint) kmlFolder() kml.Element {
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
					kml.Tessellate(true),
					kml.Width(3),
				),
			),
		)
	}
	blockBearing := -1
	switch tp.pass {
	case "N":
		blockBearing = 180
	case "S":
		blockBearing = 0
	case "E":
		blockBearing = 270
	case "W":
		blockBearing = 90
	}
	var blockPlacemark kml.Element
	if blockBearing != -1 {
		blockPlacemark = kml.Folder(
			kml.Placemark(
				kml.LineString(
					kml.Coordinates(
						center,
						sphere.FAI.Offset(center, 25000, float64(blockBearing)),
					),
				),
				kml.Style(
					kml.LineStyle(
						kml.Color(color.RGBA{R: 192, G: 0, B: 0, A: 192}),
						kml.Tessellate(true),
						kml.Width(3),
					),
				),
			),
		)
	}
	var snippet kml.Element
	switch {
	case tp.signboard:
		snippet = kml.Snippet("signboard")
	case tp.notes != "":
		snippet = kml.Snippet(tp.notes)
	case tp.pass != "":
		snippet = kml.Snippet(fmt.Sprintf("pass %s", tp.pass))
	case tp.radius != 0:
		snippet = kml.Snippet(fmt.Sprintf("%dm radius", tp.radius))
	}
	var iconStyle kml.Element
	switch {
	case tp.paddle != "":
		iconStyle = icon.PaddleIconStyle(tp.paddle)
	default:
		iconStyle = kml.IconStyle(
			kml.Icon(
				kml.Href(icon.NoneHref()),
			),
		)
	}
	return kml.Folder(
		kml.Name(tp.name),
		snippet,
		kml.Placemark(
			kml.Point(
				kml.Coordinates(center),
			),
			kml.Style(
				iconStyle,
			),
		),
		radiusPlacemark,
		blockPlacemark,
		kml.Style(
			kml.ListStyle(
				kml.ListItemType("checkHideChildren"),
			),
		),
	)
}

func (r race) kmlRouteFolder() kml.Element {
	var coordinates []kml.Coordinate
	for _, tp := range r.turnpoints {
		if !tp.offRoute {
			coordinates = append(coordinates, kml.Coordinate{Lon: tp.lon, Lat: tp.lat})
		}
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

func (r race) kmlTurnpointsFolder() kml.Element {
	var turnpointFolders []kml.Element
	for _, tp := range r.turnpoints {
		turnpointFolders = append(turnpointFolders, tp.kmlFolder())
	}
	return kml.Folder(append([]kml.Element{
		kml.Name("Turnpoints"),
	}, turnpointFolders...)...,
	)
}

func (r race) kmlDocument() kml.Element {
	return kml.KML(
		kml.Document(
			kml.Name(fmt.Sprintf("%s Route", r.name)),
			kml.Snippet(r.snippet),
			kml.Open(true),
			r.kmlRouteFolder(),
			r.kmlTurnpointsFolder(),
		),
	)
}

func (r race) xcPlannerURL() *url.URL {
	var coords [][]float64
	for _, tp := range r.turnpoints {
		coords = append(coords, []float64{tp.lat, tp.lon})
	}
	vs := url.Values{}
	vs.Set("l", "free")
	vs.Set("p", string(polyline.EncodeCoords(coords)))
	vs.Set("s", strconv.Itoa(5))
	vs.Set("a", strconv.Itoa(2000))
	return &url.URL{
		Scheme:   "https",
		Host:     "xcplanner.appspot.com",
		RawQuery: vs.Encode(),
	}
}

func dms(d, m, s float64) float64 {
	return d + m/60 + s/3600
}

func run() error {
	flag.Parse()
	r, ok := races[*raceFlag]
	if !ok {
		return fmt.Errorf("unknown race: %q", *raceFlag)
	}
	switch *formatFlag {
	case "kml":
		return r.kmlDocument().WriteIndent(os.Stdout, "", "  ")
	case "xcplanner":
		_, err := os.Stdout.WriteString(r.xcPlannerURL().String() + "\n")
		return err
	default:
		return fmt.Errorf("unknown format: %q", *formatFlag)
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
