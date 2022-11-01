package extract

import (
	"errors"
	"graph"

	"github.com/paulmach/osm"
)

var highwayTable map[string]bool

func init() {

	highwayTable = map[string]bool{
		"motorway":       true,
		"motorway_link":  true,
		"trunk":          true,
		"trunk_link":     true,
		"primary":        true,
		"primary_link":   true,
		"secondary":      true,
		"secondary_link": true,
		"tertiary":       true,
		"tertiary_link":  true,
		"unclassified":   true,
		"residential":    true,
		"living_street":  true,
		"service":        true,
		"pedestrian":     true,
	}
}

func FindTag(tags *osm.Tags, key string) (string, error) {
	for _, tag := range *tags {
		if tag.Key == key {
			return tag.Value, nil
		}
	}

	return "", errors.New("No key")
}

func ParseOSMNode(_node osm.Object) *graph.ExternalNode {
	node := _node.(*osm.Node)

	ret := &graph.ExternalNode{}

	ret.Id = int64(node.ID)
	ret.X = node.Lon
	ret.Y = node.Lat

	// ... other attribution
	return ret
}

func ParseOSMWay(_way osm.Object) *graph.ResultWay {
	way := _way.(*osm.Way)
	ret := &graph.ResultWay{}

	highway, _ := FindTag(&way.Tags, "highway")
	if highway == "" {
		return nil
	}

	if _, exist := highwayTable[highway]; !exist {
		return nil
	}

	ret.Id = int64(way.ID)
	ret.Nodes = make([]int64, 0, len(way.Nodes))
	for _, node := range way.Nodes {
		ret.Nodes = append(ret.Nodes, int64(node.ID))
	}

	oneway, _ := FindTag(&way.Tags, "oneway")
	if oneway == "yes" {
		ret.Oneway = true
	} else {
		ret.Oneway = false
	}

	// ... other attribution
	return ret
}

// func ParseOSMRestriction(_restriction osm.Object) *ParsingRestriction {
// 	restriction := _restriction.(*osm.Relation)

// 	ret := &ParsingRestriction{}
// 	ret.VIAS = make([]int, 0, 1)

// 	for _, member := range restriction.Members {
// 		switch member.Role {
// 		case "from":
// 			if member.Type != "way" {
// 				panic("from type is not way")
// 			}
// 			ret.FROM = int(member.Ref)
// 		case "via":
// 			ret.VIAS = append(ret.VIAS, int(member.Ref))
// 		case "to":
// 			if member.Type != "way" {
// 				panic("to type is not way")
// 			}
// 			ret.TO = int(member.Ref)
// 		}
// 	}

// 	if len(ret.VIAS) >= 2 {
// 		ret.MULTIRESTRICTION = 1
// 	} else if len(ret.VIAS) == 1 {
// 		ret.MULTIRESTRICTION = 0
// 	} else {
// 		ret = nil
// 	}

// 	if ret.FROM == 0 || ret.TO == 0 {
// 		ret = nil
// 	}

// 	// if ret is nil, invalid restriction
// 	return ret
// }
