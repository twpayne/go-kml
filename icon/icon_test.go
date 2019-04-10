package icon

import "testing"

func TestHref(t *testing.T) {
	for _, tc := range []struct {
		got  string
		want string
	}{
		{
			got:  CharacterHref('9'),
			want: "https://maps.google.com/mapfiles/kml/pal3/icon8.png",
		},
		{
			got:  CharacterHref('A'),
			want: "https://maps.google.com/mapfiles/kml/pal5/icon48.png",
		},
		{
			got:  CharacterHref('M'),
			want: "https://maps.google.com/mapfiles/kml/pal5/icon36.png",
		},
		{
			got:  CharacterHref('Z'),
			want: "https://maps.google.com/mapfiles/kml/pal5/icon1.png",
		},
		{
			got:  DefaultHref(),
			want: "https://maps.google.com/mapfiles/kml/pushpin/ylw-pushpin.png",
		},
		{
			got:  NoneHref(),
			want: "https://maps.google.com/mapfiles/kml/pal2/icon15.png",
		},
		{
			got:  NumberHref(1),
			want: "https://maps.google.com/mapfiles/kml/pal3/icon0.png",
		},
		{
			got:  NumberHref(10),
			want: "https://maps.google.com/mapfiles/kml/pal3/icon17.png",
		},
		{
			got:  PaddleHref("A"),
			want: "https://maps.google.com/mapfiles/kml/paddle/A.png",
		},
		{
			got:  TrackHref(0),
			want: "https://earth.google.com/images/kml-icons/track-directional/track-0.png",
		},
	} {
		if tc.got != tc.want {
			t.Errorf("got %s, want %s", tc.got, tc.want)
		}
	}
}
