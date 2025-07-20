package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"

	"github.com/twpayne/go-waypoint"

	"github.com/twpayne/go-kml/v3"
	"github.com/twpayne/go-kml/v3/icon"
)

var name = flag.String("name", "Route", "name")

func readWaypoints(filename string) (waypoint.Collection, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	wc, _, err := waypoint.Read(f)
	return wc, err
}

func run() error {
	flag.Parse()

	if flag.NArg() < 1 {
		return fmt.Errorf("syntax: %s waypoint-file [waypoints...]", os.Args[0])
	}

	waypoints, err := readWaypoints(flag.Arg(0))
	if err != nil {
		return err
	}

	waypointsByID := make(map[string]*waypoint.T)
	for _, w := range waypoints {
		if _, ok := waypointsByID[w.ID]; ok {
			return fmt.Errorf("duplicate waypoint ID: %s", w.ID)
		}
		waypointsByID[w.ID] = w
	}

	turnpoints := make([]*waypoint.T, flag.NArg()-1)
	turnpointIDs := make(map[string]bool)
	for i, arg := range flag.Args()[1:] {
		turnpoint, ok := waypointsByID[arg]
		if !ok {
			return fmt.Errorf("unknown waypoint: %s", arg)
		}
		turnpoints[i] = turnpoint
		turnpointIDs[arg] = true
	}

	routeCoordinates := make([]kml.Coordinate, len(turnpoints))
	for i, turnpoint := range turnpoints {
		coordinate := kml.Coordinate{
			Lon: turnpoint.Longitude,
			Lat: turnpoint.Latitude,
			Alt: turnpoint.Altitude,
		}
		routeCoordinates[i] = coordinate
	}

	routeFolder := kml.Folder(
		kml.Name("Route"),
		kml.Placemark(
			kml.LineString(
				kml.Coordinates(routeCoordinates...),
				kml.Tessellate(true),
			),
			kml.Style(
				kml.LineStyle(
					kml.Color(color.RGBA{R: 192, G: 0, B: 0, A: 192}),
					kml.Width(3),
				),
			),
		),
		kml.Style(
			kml.ListStyle(
				kml.ListItemType(kml.ListItemTypeCheckHideChildren),
			),
		),
	)

	turnpointFolders := make([]kml.Element, len(turnpoints))
	for i := range turnpoints {
		turnpoint := turnpoints[i]
		var name string
		var paddleID string
		if i == 0 {
			name = "START"
			paddleID = "go"
		} else {
			name = fmt.Sprintf("TP%02d", i)
			paddleID = string([]byte{byte('A' + i - 1)})
		}
		turnpointFolder := kml.Folder(
			kml.Name(fmt.Sprintf("%s %s", name, turnpoint.Description)),
			kml.Placemark(
				kml.Point(
					kml.Coordinates(kml.Coordinate{
						Lon: turnpoint.Longitude,
						Lat: turnpoint.Latitude,
						Alt: turnpoint.Altitude,
					}),
				),
				kml.Style(
					kml.IconStyle(
						kml.HotSpot(kml.Vec2{X: 0.5, Y: 0, XUnits: kml.UnitsFraction, YUnits: kml.UnitsFraction}),
						kml.Icon(
							kml.Href(icon.PaddleHref(paddleID)),
						),
						kml.Scale(0.5),
					),
				),
			),
			kml.Style(
				kml.ListStyle(
					kml.ListItemType(kml.ListItemTypeCheckHideChildren),
				),
			),
		)
		turnpointFolders[i] = turnpointFolder
	}

	turnpointsFolder := kml.Folder(
		append([]kml.Element{
			kml.Name("Turnpoints"),
			kml.Open(true),
		},
			turnpointFolders...,
		)...,
	)

	waypointFolders := make([]kml.Element, 0, len(waypoints))
	for _, waypoint := range waypoints {
		if _, ok := turnpointIDs[waypoint.ID]; ok {
			continue
		}
		waypointFolder := kml.Folder(
			kml.Name(waypoint.Description),
			kml.Placemark(
				kml.Point(
					kml.Coordinates(kml.Coordinate{
						Lon: waypoint.Longitude,
						Lat: waypoint.Latitude,
						Alt: waypoint.Latitude,
					}),
				),
				kml.Style(
					kml.IconStyle(
						kml.Icon(
							kml.Href(
								icon.PaletteHref(2, 13),
							),
						),
						kml.Scale(0.5),
					),
				),
			),
		)
		waypointFolders = append(waypointFolders, waypointFolder)
	}

	waypointsFolder := kml.Folder(
		append([]kml.Element{
			kml.Name("Waypoints"),
			kml.Open(false),
		},
			waypointFolders...,
		)...,
	)

	result := kml.KML(
		kml.Document(
			kml.Name(*name),
			kml.Open(true),
			routeFolder,
			turnpointsFolder,
			waypointsFolder,
		),
	)

	return result.WriteIndent(os.Stdout, "", "  ")
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
