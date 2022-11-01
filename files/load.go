package files

import (
	"encoding/binary"
	"graph"
	"os"
)

func LoadGeoNodes(filepath string) []graph.ExternalNode {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var count int32
	err = binary.Read(f, binary.LittleEndian, &count)

	if count == 0 {
		panic("Data count is zero")
	}

	ret := make([]graph.ExternalNode, count)
	binary.Read(f, binary.LittleEndian, ret)
	return ret
}

func LoadUncompEdges(filepath string) []graph.InternalEdge {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var count int32
	err = binary.Read(f, binary.LittleEndian, &count)
	if count == 0 {
		panic("Data count is zero")
	}

	ret := make([]graph.InternalEdge, count)
	binary.Read(f, binary.LittleEndian, ret)
	return ret
}

func LoadEdgeAnnotations(filepath string) []graph.EdgeAnnotation {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var count int32
	err = binary.Read(f, binary.LittleEndian, &count)
	if count == 0 {
		panic("Data count is zero")
	}

	ret := make([]graph.EdgeAnnotation, count)
	binary.Read(f, binary.LittleEndian, ret)
	return ret
}
