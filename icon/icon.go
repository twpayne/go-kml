// Package icon provides helper functions for standard Google Earth icons.
// See http://kml4earth.appspot.com/icons.html.
package icon

import (
	"strconv"

	"github.com/twpayne/go-kml"
)

// CharacterHref returns the href of the icon wirh the specificed character, or
// an empty string if no such icon exists. See
// http://kml4earth.appspot.com/icons.html#pal3 and
// http://kml4earth.appspot.com/icons.html#pal5.
func CharacterHref(c rune) string {
	switch {
	case '1' <= c && c <= '9':
		return PaletteHref(3, int(c-'1'))
	case 'A' <= c && c <= 'Z':
		return PaletteHref(5, int((c-'A')%8+16*((31-c+'A')/8)))
	default:
		return ""
	}
}

// DefaultHref returns the href of the default icon.
func DefaultHref() string {
	return PushpinHref("ylw")
}

// NoneHref returns the icon of the empty icon.
func NoneHref() string {
	return PaletteHref(2, 15)
}

// NumberHref returns the href of the icon with number n. See
// http://kml4earth.appspot.com/icons.html#pal3.
func NumberHref(n int) string {
	if 1 <= n && n <= 10 {
		return PaletteHref(3, (n-1)%8+16*((n-1)/8))
	}
	return ""
}

// PaddleHref returns the href of the paddle icon with id. See
// http://kml4earth.appspot.com/icons.html#paddle.
func PaddleHref(id string) string {
	return "https://maps.google.com/mapfiles/kml/paddle/" + id + ".png"
}

// PaddleIconStyle returns an IconStyle for the paddle icon with id and the
// hotspot set. See http://kml4earth.appspot.com/icons.html#paddle.
func PaddleIconStyle(id string) kml.Element {
	return kml.IconStyle(
		kml.HotSpot(kml.Vec2{X: 0.5, Y: 0, XUnits: "fraction", YUnits: "fraction"}),
		kml.Icon(
			kml.Href(PaddleHref(id)),
		),
	)
}

// PaletteHref returns the href of icon in pal.
func PaletteHref(pal, icon int) string {
	return "https://maps.google.com/mapfiles/kml/pal" + strconv.Itoa(pal) + "/icon" + strconv.Itoa(icon) + ".png"
}

// PushpinHref returns the href of pushpin of color. Valid colors are blue,
// green, ltblu, pink, purple, red, wht, and ylw. See
// http://kml4earth.appspot.com/icons.html#pushpin.
func PushpinHref(color string) string {
	return "https://maps.google.com/mapfiles/kml/pushpin/" + color + "-pushpin.png"
}

// ShapeHref returns the href of the icon with the specified shape. See
// http://kml4earth.appspot.com/icons.html#shapes.
func ShapeHref(shape string) string {
	return "http://maps.google.com/mapfiles/kml/shapes/" + shape + ".png"
}

// TrackHref returns the href of the ith track icon. See
// http://kml4earth.appspot.com/icons.html#kml-icons.
func TrackHref(i int) string {
	return "https://earth.google.com/images/kml-icons/track-directional/track-" + strconv.Itoa(i) + ".png"
}

// TrackNoneHref returns the href of the track icon when there is no heading.
// See http://kml4earth.appspot.com/icons.html#kml-icons.
func TrackNoneHref() string {
	return "https://earth.google.com/images/kml-icons/track-directional/track-none.png"
}
