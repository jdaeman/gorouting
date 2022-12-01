package engine

import (
	"errors"
	"fmt"
	"modules/files"
)

type RouteParameter struct {
	Coords [][2]float64 // [x,y];[x,y]...
}

func toNodeList(phatom PhantomNode) []int32 {
	ret := make([]int32, 0, 2)

	if phatom.IsValidForward() {
		ret = append(ret, phatom.ForwardId)
	}
	if phatom.IsValidBackward() {
		ret = append(ret, phatom.BackwardId)
	}

	return ret
}

func Route(repository *files.DataRepository, engine *RoutingEngine, params RouteParameter, result *map[string]interface{}) error {
	// parameter check
	if params.Coords == nil || len(params.Coords) == 0 {
		return errors.New("requested coordinate is empty.")
	}
	if len(params.Coords)%2 == 1 {
		return errors.New("invalid coordinates count.")
	}
	if len(params.Coords) > 2 {
		return errors.New("currently, via route is not support.")
	}

	// map matching
	phantomNodes := engine.rtree.GetPhantomNodes(params.Coords)
	if len(phantomNodes) != len(params.Coords) {
		errMsg := fmt.Sprint("Map matching error coordinate index ", len(phantomNodes))
		return errors.New(errMsg)
	}

	routePath := make([]int32, 0)

	// route
	for i := 1; i < len(phantomNodes); i++ {
		from, to := phantomNodes[i-1], phantomNodes[i]
		sources, goals := toNodeList(from), toNodeList(to)
		subPath := engine.ShortestPathSearch(sources, goals)
		if subPath == nil {
			routePath = nil
			break
		}
		routePath = append(routePath, subPath...)
	}

	if routePath == nil {
		return errors.New("Could not find path.")
	}

	// prepare data
	*result = MakeRouteResponse(repository, phantomNodes, routePath)
	return nil
}
