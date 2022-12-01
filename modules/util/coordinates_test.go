package util_test

import (
	"fmt"
	"modules/util"
	"testing"
)

func TestBearing(t *testing.T) {
	p1 := [2]float64{129.0923865, 35.2922260}
	p2 := [2]float64{129.0940302, 35.2914267}

	bearing := util.GetBearing(p1, p2)
	fmt.Println("bearing1", bearing)

	if bearing < 90 || bearing > 180 {
		t.Fail()
	}

	bearing = util.GetBearing(p2, p1)
	fmt.Println("bearing2", bearing)
	if bearing > 0 && bearing < 270 {
		t.Fail()
	}
}

func TestHaversineDist(t *testing.T) {
	p1 := [2]float64{129.0923865, 35.2922260}
	p2 := [2]float64{129.0940302, 35.2914267}

	distance1 := int(util.HaversineDistance(p1, p2))
	distance2 := int(util.HaversineDistance(p2, p1))
	diff := distance1 - distance2

	fmt.Println("dist1", distance1)
	fmt.Println("dist2", distance2)

	if diff > 1 || diff < -1 {
		t.Fail()
	}
}
