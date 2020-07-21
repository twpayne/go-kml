package icon

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHref(t *testing.T) {
	for i, tc := range []struct {
		actual   string
		expected string
	}{
		{
			actual:   CharacterHref('9'),
			expected: "https://maps.google.com/mapfiles/kml/pal3/icon8.png",
		},
		{
			actual:   CharacterHref('A'),
			expected: "https://maps.google.com/mapfiles/kml/pal5/icon48.png",
		},
		{
			actual:   CharacterHref('M'),
			expected: "https://maps.google.com/mapfiles/kml/pal5/icon36.png",
		},
		{
			actual:   CharacterHref('Z'),
			expected: "https://maps.google.com/mapfiles/kml/pal5/icon1.png",
		},
		{
			actual:   DefaultHref(),
			expected: "https://maps.google.com/mapfiles/kml/pushpin/ylw-pushpin.png",
		},
		{
			actual:   NoneHref(),
			expected: "https://maps.google.com/mapfiles/kml/pal2/icon15.png",
		},
		{
			actual:   NumberHref(1),
			expected: "https://maps.google.com/mapfiles/kml/pal3/icon0.png",
		},
		{
			actual:   NumberHref(10),
			expected: "https://maps.google.com/mapfiles/kml/pal3/icon17.png",
		},
		{
			actual:   PaddleHref("A"),
			expected: "https://maps.google.com/mapfiles/kml/paddle/A.png",
		},
		{
			actual:   TrackHref(0),
			expected: "https://earth.google.com/images/kml-icons/track-directional/track-0.png",
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.actual)
		})
	}
}
