package sphere_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-kml/v3"
	"github.com/twpayne/go-kml/v3/sphere"
)

func TestCircle(t *testing.T) {
	for i, tc := range []struct {
		center   kml.Coordinate
		radius   float64
		maxErr   float64
		expected []kml.Coordinate
	}{
		{
			center: kml.Coordinate{Lon: 0, Lat: 0, Alt: 100},
			radius: 1000,
			maxErr: 1,
			expected: []kml.Coordinate{
				{Lon: 0, Lat: 0.008983152841195214, Alt: 100},
				{Lon: 0.0011258876022698656, Lat: 0.008912317997331125, Alt: 100},
				{Lon: 0.002234019283634706, Lat: 0.008700930573621838, Alt: 100},
				{Lon: 0.0033069191447876595, Lat: 0.008352324276246322, Alt: 100},
				{Lon: 0.004327666913493653, Lat: 0.007871996835109306, Alt: 100},
				{Lon: 0.005280164787461206, Lat: 0.007267523301307612, Alt: 100},
				{Lon: 0.006149391306330964, Lat: 0.0065484365839908, Alt: 100},
				{Lon: 0.006921638249063812, Lat: 0.005726077110629261, Alt: 100},
				{Lon: 0.007584726820717381, Lat: 0.004813413981645677, Alt: 100},
				{Lon: 0.008128199719224507, Lat: 0.003824840439915492, Alt: 100},
				{Lon: 0.00854348605317887, Lat: 0.0027759468807100067, Alt: 100},
				{Lon: 0.008824036509791956, Lat: 0.001683274981853487, Alt: 100},
				{Lon: 0.008965426641358826, Lat: 0.0005640568316080706, Alt: 100},
				{Lon: 0.008965426641358826, Lat: -0.0005640568316080696, Alt: 100},
				{Lon: 0.008824036509791956, Lat: -0.0016832749818534878, Alt: 100},
				{Lon: 0.00854348605317887, Lat: -0.002775946880710006, Alt: 100},
				{Lon: 0.008128199719224507, Lat: -0.003824840439915493, Alt: 100},
				{Lon: 0.007584726820717381, Lat: -0.004813413981645679, Alt: 100},
				{Lon: 0.006921638249063812, Lat: -0.00572607711062926, Alt: 100},
				{Lon: 0.006149391306330961, Lat: -0.006548436583990801, Alt: 100},
				{Lon: 0.005280164787461206, Lat: -0.007267523301307612, Alt: 100},
				{Lon: 0.004327666913493655, Lat: -0.007871996835109304, Alt: 100},
				{Lon: 0.0033069191447876573, Lat: -0.008352324276246324, Alt: 100},
				{Lon: 0.002234019283634706, Lat: -0.008700930573621838, Alt: 100},
				{Lon: 0.001125887602269864, Lat: -0.008912317997331125, Alt: 100},
				{Lon: 1.1001189463363886e-18, Lat: -0.008983152841195214, Alt: 100},
				{Lon: -0.0011258876022698621, Lat: -0.008912317997331125, Alt: 100},
				{Lon: -0.002234019283634708, Lat: -0.008700930573621838, Alt: 100},
				{Lon: -0.003306919144787659, Lat: -0.008352324276246324, Alt: 100},
				{Lon: -0.004327666913493654, Lat: -0.007871996835109304, Alt: 100},
				{Lon: -0.005280164787461205, Lat: -0.007267523301307612, Alt: 100},
				{Lon: -0.00614939130633096, Lat: -0.006548436583990801, Alt: 100},
				{Lon: -0.006921638249063812, Lat: -0.00572607711062926, Alt: 100},
				{Lon: -0.007584726820717378, Lat: -0.004813413981645681, Alt: 100},
				{Lon: -0.00812819971922451, Lat: -0.0038248404399154876, Alt: 100},
				{Lon: -0.00854348605317887, Lat: -0.002775946880710008, Alt: 100},
				{Lon: -0.008824036509791956, Lat: -0.001683274981853488, Alt: 100},
				{Lon: -0.008965426641358826, Lat: -0.0005640568316080777, Alt: 100},
				{Lon: -0.008965426641358826, Lat: 0.0005640568316080745, Alt: 100},
				{Lon: -0.008824036509791954, Lat: 0.0016832749818534924, Alt: 100},
				{Lon: -0.00854348605317887, Lat: 0.002775946880710005, Alt: 100},
				{Lon: -0.008128199719224507, Lat: 0.003824840439915492, Alt: 100},
				{Lon: -0.0075847268207173855, Lat: 0.004813413981645672, Alt: 100},
				{Lon: -0.006921638249063809, Lat: 0.005726077110629263, Alt: 100},
				{Lon: -0.0061493913063309594, Lat: 0.006548436583990802, Alt: 100},
				{Lon: -0.005280164787461206, Lat: 0.007267523301307612, Alt: 100},
				{Lon: -0.004327666913493653, Lat: 0.007871996835109306, Alt: 100},
				{Lon: -0.0033069191447876655, Lat: 0.008352324276246322, Alt: 100},
				{Lon: -0.002234019283634703, Lat: 0.00870093057362184, Alt: 100},
				{Lon: -0.0011258876022698613, Lat: 0.008912317997331125, Alt: 100},
				{Lon: 0, Lat: 0.008983152841195214, Alt: 100},
			},
		},
		{
			center: kml.Coordinate{Lon: 13.631333, Lat: 46.438500},
			radius: 50,
			maxErr: 1,
			expected: []kml.Coordinate{
				{Lon: 13.631333, Lat: 46.43894915764205},
				{Lon: 13.631658888465811, Lat: 46.43888898146551},
				{Lon: 13.631897453677293, Lat: 46.43872457743259},
				{Lon: 13.631984772278717, Lat: 46.438499998148764},
				{Lon: 13.631897449024441, Lat: 46.43827541979054},
				{Lon: 13.63165888381296, Lat: 46.43811101760886},
				{Lon: 13.631333, Lat: 46.43805084235793},
				{Lon: 13.63100711618704, Lat: 46.43811101760886},
				{Lon: 13.630768550975558, Lat: 46.43827541979054},
				{Lon: 13.630681227721285, Lat: 46.438499998148764},
				{Lon: 13.630768546322708, Lat: 46.43872457743259},
				{Lon: 13.63100711153419, Lat: 46.43888898146551},
				{Lon: 13.631333, Lat: 46.43894915764205},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := sphere.WGS84.Circle(tc.center, tc.radius, tc.maxErr)
			assert.Equal(t, len(tc.expected), len(actual))
			for i, actualCoordinate := range actual {
				assertInDelta(t, tc.expected[i].Lon, actualCoordinate.Lon, 1e-14)
				assertInDelta(t, tc.expected[i].Lat, actualCoordinate.Lat, 1e-14)
				assert.Equal(t, tc.center.Alt, actualCoordinate.Alt)
				assertInDelta(t, tc.radius, sphere.WGS84.HaversineDistance(tc.center, actualCoordinate), 1e-9)
			}
			for _, expectedCoordinate := range tc.expected {
				assertInDelta(t, tc.radius, sphere.WGS84.HaversineDistance(tc.center, expectedCoordinate), 1e-9)
			}
		})
	}
}

func assertInDelta(tb testing.TB, expected, actual, delta float64) {
	tb.Helper()
	if math.Abs(expected-actual) <= delta {
		return
	}
	tb.Fatalf("Expected %f to be within %f of %f", actual, delta, expected)
}
