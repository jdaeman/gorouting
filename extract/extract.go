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

	nodeBasedGraphFactory := NewNodeBasedGraphFactory(config.OsmPath)
	edgeBasedGraphFactory := NewEdgeBasedGraphFactory(nodeBasedGraphFactory)
	edgeBasedGraphFactory.Run()

	// RTree
	// SCC (?)

	return 0
}

func ParseOSMData(config Config) bool {
	log.Println("Read file from", config.OsmPath)
	objects, err := ReadOSM(config.OsmPath)
	if err != nil {
		log.Println("Error", err)
		return false
	}

	nodes, ways, relations := objects[0], objects[1], objects[2]
	extractor := NewExtractor(nodes, ways, relations, config.OsmPath)
	nodes, ways, relations, objects = nil, nil, nil, nil

	// parse each objects.
	nodeCount := extractor.ProcessOSMNodes()
	log.Println("Raw osm node count", nodeCount)
	wayCount := extractor.ProcessOSMWays()
	log.Println("Raw osm drivable way count", wayCount)
	restrictionCount := extractor.ProcessOSMRestrictions()
	log.Println("Raw osm restriction count", restrictionCount)

	// internal process step1.
	extractor.ProcessNodes()
	extractor.ProcessEdges()
	extractor.ProcessRestrictions()

	// prepare internal graph data..
	// save as file.
	extractor.PrepareData()

	return true
}

// Read osm file.
// Return Node, Way, Relation object slices.
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
