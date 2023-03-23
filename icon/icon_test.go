package icon_test

import (
	"strconv"
	"testing"

	"github.com/alecthomas/assert"

	"github.com/twpayne/go-kml/v2/icon"
)

func TestHref(t *testing.T) {
	for i, tc := range []struct {
		actual   string
		expected string
	}{
		{
			actual:   icon.CharacterHref('9'),
			expected: "https://maps.google.com/mapfiles/kml/pal3/icon8.png",
		},
		{
			actual:   icon.CharacterHref('A'),
			expected: "https://maps.google.com/mapfiles/kml/pal5/icon48.png",
		},
		{
			actual:   icon.CharacterHref('M'),
			expected: "https://maps.google.com/mapfiles/kml/pal5/icon36.png",
		},
		{
			actual:   icon.CharacterHref('Z'),
			expected: "https://maps.google.com/mapfiles/kml/pal5/icon1.png",
		},
		{
			actual:   icon.DefaultHref(),
			expected: "https://maps.google.com/mapfiles/kml/pushpin/ylw-pushpin.png",
		},
		{
			actual:   icon.NoneHref(),
			expected: "https://maps.google.com/mapfiles/kml/pal2/icon15.png",
		},
		{
			actual:   icon.NumberHref(1),
			expected: "https://maps.google.com/mapfiles/kml/pal3/icon0.png",
		},
		{
			actual:   icon.NumberHref(10),
			expected: "https://maps.google.com/mapfiles/kml/pal3/icon17.png",
		},
		{
			actual:   icon.PaddleHref("A"),
			expected: "https://maps.google.com/mapfiles/kml/paddle/A.png",
		},
		{
			actual:   icon.TrackHref(0),
			expected: "https://earth.google.com/images/kml-icons/track-directional/track-0.png",
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.actual)
		})
	}
}
