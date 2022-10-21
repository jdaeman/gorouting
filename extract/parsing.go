package extract

import (
	"errors"

	"github.com/paulmach/osm"
)

type ParsingNode struct {
	ID int
	X  float64
	Y  float64
}

type ParsingWay struct {
	ID    int
	Nodes []int

	Oneway bool
	// ...
}

type ParsingRestriction struct {
	FROM int
	VIAS []int
	TO   int

	MULTIRESTRICTION int // no: 0, yes: 1
}

func FindTag(tags *osm.Tags, key string) (string, error) {
	for _, tag := range *tags {
		if tag.Key == key {
			return tag.Value, nil
		}
	}

	return "", errors.New("No key")
}

func ParseOSMNode(_node osm.Object) *ParsingNode {
	node := _node.(*osm.Node)

	ret := &ParsingNode{}

	ret.ID = int(node.ID)
	ret.X = node.Lon
	ret.Y = node.Lat

	// ... other attribution

	return ret
}

func ParseOSMWay(_way osm.Object) *ParsingWay {
	way := _way.(*osm.Way)

	ret := &ParsingWay{}

	ret.ID = int(way.ID)
	ret.Nodes = make([]int, 0, len(way.Nodes))
	for _, node := range way.Nodes {
		ret.Nodes = append(ret.Nodes, int(node.ID))
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

func ParseOSMRestriction(_restriction osm.Object) *ParsingRestriction {
	restriction := _restriction.(*osm.Relation)

	ret := &ParsingRestriction{}
	ret.VIAS = make([]int, 0, 1)

	for _, member := range restriction.Members {
		switch member.Role {
		case "from":
			if member.Type != "way" {
				panic("from type is not way")
			}
			ret.FROM = int(member.Ref)
		case "via":
			ret.VIAS = append(ret.VIAS, int(member.Ref))
		case "to":
			if member.Type != "way" {
				panic("to type is not way")
			}
			ret.TO = int(member.Ref)
		}
	}

	if len(ret.VIAS) >= 2 {
		ret.MULTIRESTRICTION = 1
	} else if len(ret.VIAS) == 1 {
		ret.MULTIRESTRICTION = 0
	} else {
		ret = nil
	}

	if ret.FROM == 0 || ret.TO == 0 {
		ret = nil
	}

	// if ret is nil, invalid restriction
	return ret
}
