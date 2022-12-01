package engine

import (
	"modules/files"
)

type maneuver_obj struct {
	Bearing1 int       `json:"bearing_after"`
	Bearing2 int       `json:"bearing_before"`
	Type     string    `json:"type"`
	Location []float64 `json:"location"`
}

type intersection_obj struct {
	Entry    []bool    `json:"entry"`
	Bearings []int     `json:"bearings"`
	Location []float64 `json:"location"`
}

type step_obj struct {
	Distance      int                `json:"distance"`
	Duration      int                `json:"duration"`
	Geometry      string             `json:"geometry"`
	Intersections []intersection_obj `json:"intersections"`
	Mode          string             `json:"mode"`
	Name          string             `json:"name"`
	Ref           string             `json:"ref"`
	Maneuver      []maneuver_obj     `json:"maneuver"`
}

type leg_obj struct {
	Distance int `json:"distance"`
	Duration int `json:"duraiton"`

	Steps   []step_obj `json:"steps"`
	Summary string     `json:"summary"`
}

type route_obj struct {
	Distance   int       `json:"distance"`
	Duration   int       `json:"duration"`
	Legs       []leg_obj `json:"legs"`
	Weight     int       `json:"weight"`
	WeightName string    `json:"weight_name"`
}

type waypoint_obj struct {
	Name     string    `json:"name"`
	Location []float64 `json:"location"`
}

func MakeRouteResponse(
	repository *files.DataRepository,
	phantoms []PhantomNode, path []int32) map[string]interface{} {

	//var defaultCodec = polyline.Codec{Dim: 3, Scale: 1e5}

	ret := make(map[string]interface{})

	wayIds := make([]int64, 0)
	geolocations := make([][]float64, 0)
	polylines := make([][]float64, 0)

	for _, u := range path {
		geoId := repository.GetGeometryId(u)
		annoId := repository.GetAnnotationId(u)
		wayId := repository.GetAnnotation(annoId).Id
		geoNodeIds := repository.GetGeoNodeIds(geoId)
		locations := repository.GetLocations(geoNodeIds)
		polyline := repository.GetLocationsYX(geoNodeIds)

		polylines = append(polylines, polyline...)
		geolocations = append(geolocations, locations...)
		wayIds = append(wayIds, wayId)
	}

	// ret["code"] = "ok"
	// ret["geometry"] = geolocations
	// ret["way"] = wayIds
	// ret["polyline"] = string(defaultCodec.EncodeCoords(nil, polylines))

	routes := make([]route_obj, 1)
	route := &routes[0]

	route.Distance = 111
	route.Duration = 222
	route.Weight = 333
	route.WeightName = "distance"
	route.Legs = make([]leg_obj, 1)

	leg := &route.Legs[0]

	leg.Distance = 11
	leg.Duration = 22
	leg.Steps = make([]step_obj, 2)

	step := &leg.Steps[0]
	step.Distance = 1
	step.Duration = 2
	step.Geometry = "iyjdFovhfW]\\tI`O`MQbCfBlJx^lHp`@pDjVHhLyCrAwC|G_D`Ab@nLeCE@"
	step.Intersections = make([]intersection_obj, 1)
	step.Maneuver = make([]maneuver_obj, 1)
	step.Mode = "driving"
	step.Name = "road"
	step.Ref = "3"

	step = &leg.Steps[1]
	step.Distance = 1
	step.Duration = 2
	step.Geometry = "iyjdFovhfW]\\tI`O`MQbCfBlJx^lHp`@pDjVHhLyCrAwC|G_D`Ab@nLeCE@"
	step.Intersections = make([]intersection_obj, 1)
	step.Maneuver = make([]maneuver_obj, 1)
	step.Mode = "driving"
	step.Name = "road"
	step.Ref = "3"

	waypoints := make([]waypoint_obj, 2)
	waypoints[0].Location = make([]float64, 2)
	waypoints[1].Location = make([]float64, 2)

	for i := 0; i < 2; i++ {
		waypoints[i].Location[0] = phantoms[i].X
		waypoints[i].Location[1] = phantoms[i].Y
	}

	ret["code"] = "Ok"
	ret["routes"] = routes
	ret["waypoints"] = waypoints

	return ret
}
