package ellipsoid

// FIXME add Vincenty algorithm https://github.com/twpayne/ol3/blob/189ba34955029280ba348a65540e5e36539477f9/src/ol/ellipsoid/ellipsoid.js#L56
// FIXME compare with https://github.com/StefanSchroeder/Golang-Ellipsoid

import (
	"math"

	"github.com/twpayne/go-kml/v3"
)

// An Ellipsoid is an ellipsoid.
type Ellipsoid struct {
	A          float64
	Flattening float64
}

// WGS84 is the WGS84 ellipsoid.
var WGS84 = Ellipsoid{
	A:          6378137,
	Flattening: 1.0 / 298.257223563,
}

// Distance returns the distance between c1 and c2.
func (e Ellipsoid) Distance(c1, c2 kml.Coordinate) float64 {
	// See "Annex C Ellipsoid distance between two points" in the "FAI Sporting
	// Code: Section 7F – XC Scoring".
	//
	// https://fai.org/sites/default/files/civl/documents/sporting_code_s7_f_-_xc_scoring_2025_v1.0.pdf
	lat1r := math.Pi * c1.Lat / 180
	lon1r := math.Pi * c1.Lon / 180
	lat2r := math.Pi * c2.Lat / 180
	lon2r := math.Pi * c2.Lon / 180
	if lon1r == lon2r && lat1r == lat2r {
		return 0
	}
	theta1 := math.Atan(e.oneMinusF() * math.Tan(lat1r))
	theta2 := math.Atan(e.oneMinusF() * math.Tan(lat2r))
	thetaM := (theta1 + theta2) / 2.0
	dThetaM := (theta2 - theta1) / 2.0
	dLambda := lon2r - lon1r
	dLambdaM := dLambda / 2.0
	sinThetaM := math.Sin(thetaM)
	cosThetaM := math.Cos(thetaM)
	sinDThetaM := math.Sin(dThetaM)
	cosDThetaM := math.Cos(dThetaM)
	sin2ThetaM := sinThetaM * sinThetaM
	cos2ThetaM := cosThetaM * cosThetaM
	sin2DThetaM := sinDThetaM * sinDThetaM
	cos2DThetaM := cosDThetaM * cosDThetaM
	sinDLambdaM := math.Sin(dLambdaM)
	sin2DLambdaM := sinDLambdaM * sinDLambdaM
	H := cos2ThetaM - sin2DThetaM
	L := sin2DThetaM + H*sin2DLambdaM
	cosD := 1.0 - 2.0*L
	d := math.Acos(cosD)
	sinD := math.Sin(d)
	oneMinusL := 1.0 - L
	if sinD == 0.0 || L == 0.0 || oneMinusL == 0.0 {
		return 0
	}
	U := 2.0 * sin2ThetaM * cos2DThetaM / oneMinusL
	V := 2.0 * sin2DThetaM * cos2ThetaM / L
	X := U + V
	Y := U - V
	T := d / sinD
	D := 4.0 * T * T
	E := 2.0 * cosD
	A := D * E
	B := 2.0 * D
	C := T - (A-E)/2.0
	n1 := X * (A + C*X)
	n2 := Y * (B + E*Y)
	n3 := D * X * Y
	delta1d := e.Flattening * (T*X - Y) / 4.0
	delta2d := e.flatSq64() * (n1 - n2 + n3)
	return e.A * sinD * (T - delta1d + delta2d)
}

func (e Ellipsoid) flatSq64() float64 {
	return e.Flattening * e.Flattening / 64
}

func (e Ellipsoid) oneMinusF() float64 {
	return 1 - e.Flattening
}
