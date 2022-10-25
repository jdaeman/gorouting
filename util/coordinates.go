package util

import "math"

const (
	EarthRadius = 6371
)

func degreeToRadian(deg float64) float64 {
	// 180' == 1 pi rad
	return deg * (math.Pi / 180.)
}

func radianToDegree(rad float64) float64 {
	return rad * (180. / math.Pi)
}

// [2]: {X, Y} = {Lon, Lat}
func HaversineDistance(p1 [2]float64, p2 [2]float64) float64 {
	// https://stackoverflow.com/questions/27928/calculate-distance-between-two-latitude-longitude-points-haversine-formula
	X1, X2 := p1[0], p2[0]
	Y2, Y1 := p1[1], p2[1]

	dLat := degreeToRadian(Y2 - Y1)
	dLon := degreeToRadian(X2 - X1)

	a := math.Sin(dLat/2.)*math.Sin(dLat/2) +
		math.Cos(degreeToRadian(Y1))*math.Cos(degreeToRadian(Y2))*
			math.Sin(dLon/2.)*math.Sin(dLon/2.)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := EarthRadius * c // unit is km

	return d * 1000.
}

// [2]: {X, Y} = {Lon, Lat}
func GetBearing(p1 [2]float64, p2 [2]float64) float64 {
	// https://www.igismap.com/formula-to-find-bearing-or-heading-angle-between-two-points-latitude-longitude/
	// https://www.movable-type.co.uk/scripts/latlong.html
	for i := range p1 {
		p1[i] = degreeToRadian(p1[i])
		p2[i] = degreeToRadian(p2[i])
	}

	X := math.Cos(p2[1]) * math.Sin(p2[0]-p1[0])
	Y := (math.Cos(p1[1]) * math.Sin(p2[1])) - (math.Sin(p1[1]) * math.Cos(p2[1]) * math.Cos(p2[0]-p1[0]))

	bearing := math.Atan2(X, Y)
	bearing = radianToDegree(bearing)

	if bearing < 0 {
		bearing += 360.
	}

	return bearing
}
