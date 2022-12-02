package engine

import (
	"modules/files"

	"github.com/twpayne/go-polyline"
)

func makeWayPoints(phantoms []PhantomNode) []respWayPoint {
	ret := make([]respWayPoint, len(phantoms))

	for i, phantom := range phantoms {
		ret[i].Location[0] = phantom.X
		ret[i].Location[1] = phantom.Y
		ret[i].Name = ""
	}

	return ret
}

func makeResponse(locations [][]float64) []respRoute {
	var defaultCodec = polyline.Codec{Dim: 3, Scale: 1e5}
	// make for encoded polyline.
	polyLine := make([][]float64, len(locations))
	for i, loc := range locations {
		polyLine[i] = make([]float64, 2)
		polyLine[i][0] = loc[1]
		polyLine[i][1] = loc[0]
	}

	routes := make([]respRoute, 1)
	route := &routes[0]

	route.Distance = 1
	route.Duration = 2
	route.Weight = 3
	route.WeightName = "distance"

	route.Legs = make([]respLeg, 1)
	leg := &route.Legs[0]

	leg.Distance = 1
	leg.Duration = 2
	leg.Weight = 3
	leg.Summary = "test"
	leg.Steps = make([]respStep, 2) // currently, fixed count is 2.

	intersectionPos := [2]int{len(polyLine) / 2, len(polyLine)}
	manType := [2]string{"depart", "arrive"}
	prevCut := 0

	for i := range leg.Steps {
		step := &leg.Steps[i]
		step.Distance = 1
		step.Duration = 2
		step.Weight = 3

		cut := intersectionPos[i]
		step.Geometry = string(defaultCodec.EncodeCoords(nil, polyLine[prevCut:cut]))
		step.Maneuver.Type = manType[i]

		prevCut = cut
	}

	return routes
}

func MakeRouteResponse(
	repository *files.DataRepository,
	phantoms []PhantomNode, path []int32, reqStep bool) map[string]interface{} {

	ret := make(map[string]interface{})

	wayIds := make([]int64, 0)
	geolocations := make([][]float64, 0)
	totalDist := 0

	for _, u := range path {
		// u is edge_based_node id.
		// collect data for result json.
		geoId := repository.GetGeometryId(u)
		annoId := repository.GetAnnotationId(u)

		wayId := repository.GetAnnotation(annoId).Id
		geoNodeIds, distances := repository.GetGeometry(geoId)
		locations := repository.GetLocations(geoNodeIds)

		if u == phantoms[0].ForwardId || u == phantoms[0].BackwardId {
			pos := phantoms[0].FwdPosition
			locations = locations[pos:]
		} else if u == phantoms[1].ForwardId || u == phantoms[1].BackwardId {
			pos := phantoms[1].FwdPosition
			locations = locations[:pos]
		}

		geolocations = append(geolocations, locations...)
		wayIds = append(wayIds, wayId)

		for _, segDist := range distances {
			totalDist += int(segDist)
		}
	}

	routes := makeResponse(geolocations)
	waypoints := makeWayPoints(phantoms)

	ret["code"] = "Ok"
	ret["routes"] = routes
	ret["waypoints"] = waypoints

	return ret
}
