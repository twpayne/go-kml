package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"
	"strconv"

	"github.com/twpayne/go-kml"
	"github.com/twpayne/go-kml/icon"
	"github.com/twpayne/go-waypoint"
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
		return fmt.Errorf("syntax: %s waypoint-file waypoints...", os.Args[0])
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

	var turnpoints []*waypoint.T
	turnpointIDs := make(map[string]bool)
	for _, arg := range flag.Args()[1:] {
		turnpoint, ok := waypointsByID[arg]
		if !ok {
			return fmt.Errorf("unknown waypoint: %s", arg)
		}
		turnpoints = append(turnpoints, turnpoint)
		turnpointIDs[arg] = true
	}

	var routeCoordinates []kml.Coordinate
	for _, turnpoint := range turnpoints {
		coordinate := kml.Coordinate{
			Lon: turnpoint.Longitude,
			Lat: turnpoint.Latitude,
			Alt: turnpoint.Altitude,
		}
		routeCoordinates = append(routeCoordinates, coordinate)
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

	var turnpointFolders []kml.Element
	for i, turnpoint := range turnpoints {
		var name string
		var iconStyle kml.Element
		switch i {
		case 0:
			name = "START"
			iconStyle = icon.PaddleIconStyle("go")
		case len(turnpoints) - 1:
			name = "GOAL"
			iconStyle = icon.PaddleIconStyle("stop")
		default:
			name = fmt.Sprintf("TP%02d", i)
			iconStyle = icon.PaddleIconStyle(strconv.Itoa(i))
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
					iconStyle,
				),
			),
			kml.Style(
				kml.ListStyle(
					kml.ListItemType(kml.ListItemTypeCheckHideChildren),
				),
			),
		)
		turnpointFolders = append(turnpointFolders, turnpointFolder)
	}

	turnpointsFolder := kml.Folder(
		append([]kml.Element{
			kml.Name("Turnpoints"),
			kml.Open(true),
		},
			turnpointFolders...,
		)...,
	)

	var waypointFolders []kml.Element
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
