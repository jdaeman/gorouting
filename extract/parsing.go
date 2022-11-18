package extract

import (
	"errors"
	"graph"
	"strings"

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

func ParseOSMNode(_node osm.Object) *graph.ResultNode {
	node := _node.(*osm.Node)

	ret := &graph.ResultNode{}

	ret.Id = int64(node.ID)
	ret.X = node.Lon
	ret.Y = node.Lat

	// ... other attribution
	return ret
}

// if this way is undrivable, return nil.
func ParseOSMWay(_way osm.Object) *graph.ResultWay {
	way := _way.(*osm.Way)
	ret := &graph.ResultWay{}

	highway, _ := FindTag(&way.Tags, "highway")
	if highway == "" {
		return nil
	} else if ok, exist := highwayTable[highway]; !exist {
		return nil
	} else if !ok {
		return nil
	}

	ret.Id = int64(way.ID)
	ret.Nodes = make([]int64, 0, len(way.Nodes))
	for _, node := range way.Nodes {
		ret.Nodes = append(ret.Nodes, int64(node.ID))
	}
	if len(way.Nodes) <= 1 {
		return nil
	}

	oneway, _ := FindTag(&way.Tags, "oneway")
	if oneway == "yes" {
		ret.Oneway = true
	} else {
		// bidirection.
		ret.Oneway = false
	}

	// ... other attribution
	return ret
}

// current support node restriction.
// TBD. way restriction.
func ParseOSMRelation(relation osm.Object) *graph.ResultRestriction {
	restriction := relation.(*osm.Relation)
	restriction_type, _ := FindTag(&restriction.Tags, "restriction")

	onlyRestriction := false
	if strings.Contains(restriction_type, "only") {
		onlyRestriction = true
	}

	var from, via, to int64
	from, via, to = -1, -1, -1

	for _, member := range restriction.Members {
		switch member.Role {
		case "from":
			// way
			if member.Type != "way" {
				panic("from type is not way")
			}
			from = int64(member.Ref)
		case "via":
			if member.Type == "way" {
				break
			}
			via = int64(member.Ref)
		case "to":
			if member.Type != "way" {
				panic("to type is not way")
			}
			to = int64(member.Ref)
		}
	}

	if from == -1 || via == -1 || to == -1 {
		return nil
	}

	ret := &graph.ResultRestriction{From: from, Via: via, To: to, Only: onlyRestriction}
	return ret
}
