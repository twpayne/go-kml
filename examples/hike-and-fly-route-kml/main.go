// hike-and-fly-route prints a KML file of the route of popular races.
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"image/color"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/twpayne/go-gpx"
	kml "github.com/twpayne/go-kml/v2"
	"github.com/twpayne/go-kml/v2/icon"
	"github.com/twpayne/go-kml/v2/sphere"
	polyline "github.com/twpayne/go-polyline"
)

var (
	formatFlag = flag.String("format", "kml", "format")
	raceFlag   = flag.String("race", "red-bull-x-alps-2021", "race")
)

var blockBearings = map[string]int{
	"S":  0,
	"SW": 45,
	"W":  90,
	"NW": 135,
	"N":  180,
	"NE": 225,
	"E":  270,
	"SE": 315,
}

type turnpoint struct {
	name      string
	lat       float64
	lon       float64
	ele       float64
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
	"red-bull-x-alps-2021": {
		name:    "Red Bull X-Alps 2021",
		snippet: "Created by twpayne@gmail.com",
		turnpoints: []turnpoint{
			{
				name:   "Mozartplatz",
				lat:    47.798873,
				lon:    13.047720,
				ele:    432,
				paddle: "go",
			},
			{
				name:      "Gaisberg",
				lat:       47.804398,
				lon:       13.110690,
				ele:       1275,
				paddle:    "1",
				signboard: true,
			},
			{
				name:      "Kleinarl Fußballplatz",
				lat:       47.274628,
				lon:       13.318581,
				ele:       1009,
				paddle:    "2",
				signboard: true,
			},
			{
				name:      "Kitzbühl Streif Mausefalle",
				lat:       47.426461,
				lon:       12.371147,
				ele:       1633,
				paddle:    "3",
				signboard: true,
			},
			{
				name:   "Chiemsee",
				lat:    47.858077,
				lon:    12.500269,
				ele:    521,
				radius: 3000,
				paddle: "4",
			},
			{
				name:      "Marquartstein",
				lat:       47.767503,
				lon:       12.457437,
				ele:       542,
				paddle:    "red-circle",
				signboard: true,
			},
			{
				name:     "Zugspitze",
				lat:      47.421063,
				lon:      10.985517,
				ele:      2873,
				offRoute: true,
				pass:     "N",
			},
			{
				name:      "Lermoos",
				lat:       47.401283,
				lon:       10.879767,
				ele:       990,
				paddle:    "5",
				signboard: true,
			},
			{
				name:   "Säntis",
				lat:    47.249365,
				lon:    9.343238,
				ele:    2500,
				paddle: "6",
				radius: 2000,
			},
			{
				name:      "Fiesch",
				lat:       46.40940,
				lon:       8.13688,
				ele:       1057,
				paddle:    "7",
				signboard: true,
			},
			{
				name:   "Dent d’Oche",
				lat:    46.352357,
				lon:    6.731626,
				ele:    2079,
				paddle: "8",
				pass:   "NW",
			},
			{
				name:   "Mont Blanc",
				lat:    45.830359,
				lon:    6.867674,
				ele:    4714,
				paddle: "9",
				pass:   "SW",
			},
			{
				name:   "Piz Palü",
				lat:    46.378200,
				lon:    9.958730,
				ele:    3901,
				paddle: "A",
				radius: 3500,
			},
			{
				name:      "Kronplatz",
				lat:       46.737598,
				lon:       11.954900,
				ele:       2258,
				paddle:    "B",
				signboard: true,
			},
			{
				name:      "Schmittenhöhe",
				lat:       47.328744,
				lon:       12.737518,
				ele:       1950,
				paddle:    "C",
				signboard: true,
			},
			{
				name:   "Zell am See",
				lat:    47.325290,
				lon:    12.801694,
				ele:    751,
				paddle: "stop",
			},
		},
	},
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
				paddle: "1",
			},
			{
				name:   "Tierberglihütte",
				lat:    46.702018,
				lon:    8.41421,
				paddle: "2",
			},
			{
				name:   "Lobhornhütte",
				lat:    lobhornhutteNW.Lat,
				lon:    lobhornhutteNW.Lon,
				paddle: "3",
			},
			{
				name:   "Niesen",
				lat:    46.644999,
				lon:    7.651387,
				paddle: "4",
			},
			{
				name:   "Doldehornhütte",
				lat:    46.486806,
				lon:    7.697366,
				paddle: "5",
			},
			{
				name:   "Schmadrihütte",
				lat:    46.499159,
				lon:    7.892225,
				paddle: "6",
			},
			{
				name:   "Berghaus Bäregg",
				lat:    berghausBareggSE.Lat,
				lon:    berghausBareggSE.Lon,
				paddle: "7",
			},
			{
				name:   "Glecksteinhütte",
				lat:    46.625129,
				lon:    8.096503,
				paddle: "8",
			},
			{
				name:   "Lobhornhütte",
				lat:    lobhornhutteSE.Lat,
				lon:    lobhornhutteSE.Lon,
				paddle: "9",
			},
			{
				name:   "Berghaus Bäregg",
				lat:    berghausBareggNW.Lat,
				lon:    berghausBareggNW.Lon,
				paddle: "10",
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

func (tp turnpoint) desc() string {
	switch {
	case tp.signboard:
		return "signboard"
	case tp.notes != "":
		return tp.notes
	case tp.pass != "":
		return fmt.Sprintf("pass %s", tp.pass)
	case tp.radius != 0:
		return fmt.Sprintf("%dm radius", tp.radius)
	default:
		return ""
	}
}

func (tp turnpoint) kmlFolder() kml.Element {
	center := kml.Coordinate{Lon: tp.lon, Lat: tp.lat}
	var radiusPlacemark kml.Element
	if tp.radius != 0 {
		radiusPlacemark = kml.Placemark(
			kml.LineString(
				kml.Coordinates(sphere.FAI.Circle(center, float64(tp.radius), 1)...),
				kml.Tessellate(true),
			),
			kml.Style(
				kml.LineStyle(
					kml.Color(color.RGBA{R: 0, G: 192, B: 0, A: 192}),
					kml.Width(3),
				),
			),
		)
	}
	var blockPlacemark kml.Element
	if blockBearing, ok := blockBearings[tp.pass]; ok {
		blockPlacemark = kml.Folder(
			kml.Placemark(
				kml.LineString(
					kml.Coordinates(
						center,
						sphere.FAI.Offset(center, 25000, float64(blockBearing)),
					),
					kml.Tessellate(true),
				),
				kml.Style(
					kml.LineStyle(
						kml.Color(color.RGBA{R: 192, G: 0, B: 0, A: 192}),
						kml.Width(3),
					),
				),
			),
		)
	}
	var snippet kml.Element
	if desc := tp.desc(); desc != "" {
		snippet = kml.Snippet(desc)
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
				kml.ListItemType(kml.ListItemTypeCheckHideChildren),
			),
		),
	)
}

func (r race) gpx() *gpx.GPX {
	var wpts []*gpx.WptType
	rte := &gpx.RteType{
		Name: r.name,
		Desc: r.snippet,
	}
	for _, tp := range r.turnpoints {
		if tp.offRoute {
			continue
		}
		wpt := &gpx.WptType{
			Lat:  tp.lat,
			Lon:  tp.lon,
			Ele:  tp.ele,
			Name: tp.name,
			Desc: tp.desc(),
		}
		wpts = append(wpts, wpt)
		rte.RtePt = append(rte.RtePt, wpt)
	}
	return &gpx.GPX{
		Version: "1.0",
		Creator: "ExpertGPS 1.1 - http://www.topografix.com",
		Wpt:     wpts,
		Rte:     []*gpx.RteType{rte},
	}
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
				kml.ListItemType(kml.ListItemTypeCheckHideChildren),
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
		kml.Open(true),
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
	case "gpx":
		os.Stdout.WriteString(xml.Header)
		return r.gpx().WriteIndent(os.Stdout, "", "  ")
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
