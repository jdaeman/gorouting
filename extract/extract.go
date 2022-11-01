package extract

import (
	"context"
	"log"
	"os"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmxml"
)

const (
	BUFFER_NODE     = 100000
	BUFFER_WAY      = 100000
	BUFFER_RELATION = 100000
)

func Run(config Config) int {
	ret := ParseOSMData(config)
	if !ret {
		return 1
	}

	// compress edge
	//

	return 0
}

func ParseOSMData(config Config) bool {
	log.Println("Read file from", config.OsmPath)
	objects, err := ReadOSM(config.OsmPath)

	if err != nil {
		log.Println("Error", err)
		return false
	}

	nodes := objects[0]
	ways := objects[1]
	relations := objects[2]
	extractor := NewExtractor(nodes, ways, relations, config.OsmPath)
	nodes = nil
	ways = nil
	relations = nil
	objects = nil

	nodeCount := extractor.ProcessOSMNodes()
	log.Println("Raw osm node count", nodeCount)
	wayCount := extractor.ProcessOSMWays()
	log.Println("Raw osm way count", wayCount)
	restrictionCount := extractor.ProcessOSMRestriction()
	log.Println("Raw osm restriction count", restrictionCount)

	extractor.ProcessNodes()
	extractor.ProcessEdges()
	log.Println("Used node count", len(extractor.UsedNodes))
	log.Println("Used edge count", len(extractor.UsedEdges))

	extractor.PrepareData()

	return true
}

func ReadOSM(filePath string) ([]osm.Objects, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	nodes := make(osm.Objects, 0, BUFFER_NODE)
	ways := make(osm.Objects, 0, BUFFER_WAY)
	relations := make(osm.Objects, 0, BUFFER_RELATION)

	scanner := osmxml.New(context.Background(), f)
	defer scanner.Close()

	for scanner.Scan() {
		o := scanner.Object()

		switch o.ObjectID().Type() {
		case osm.TypeNode:
			node := o.(*osm.Node)
			nodes = append(nodes, node)
		case osm.TypeWay:
			way := o.(*osm.Way)
			ways = append(ways, way)
		case osm.TypeRelation:
			relation := o.(*osm.Relation)
			rel_type, _ := FindTag(&relation.Tags, "type")
			if rel_type == "restriction" {
				relations = append(relations, relation)
			}
		}
	}

	return []osm.Objects{nodes, ways, relations}, nil
}
