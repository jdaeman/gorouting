package files

import (
	"encoding/binary"
	"modules/graph"
	"os"
)

func LoadGeoNodes(filepath string) []graph.ResultNode {
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

	ret := make([]graph.ResultNode, count)
	binary.Read(f, binary.LittleEndian, ret)
	return ret
}

func LoadEdges(filepath string) []graph.InternalEdge {
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

func LoadEdgeGeometries(filepath string) []graph.Geometry {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var count, segCount int32
	err = binary.Read(f, binary.LittleEndian, &count)
	if count == 0 {
		panic("Data count is zero")
	}

	ret := make([]graph.Geometry, count)
	for i := range ret {
		err = binary.Read(f, binary.LittleEndian, &segCount)

		if err != nil {
			panic(err)
		}

		ret[i].Nodes = make([]int32, segCount)
		ret[i].Distances = make([]int32, segCount)

		binary.Read(f, binary.LittleEndian, ret[i].Nodes)
		binary.Read(f, binary.LittleEndian, ret[i].Distances)
	}

	return ret
}

func LoadRestrictions(filepath string) []graph.InternalRestriction {
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

	ret := make([]graph.InternalRestriction, count)
	binary.Read(f, binary.LittleEndian, ret)
	return ret
}

func LoadEdgeBasedNodes(filepath string) []graph.EdgeBasedNode {
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

	ret := make([]graph.EdgeBasedNode, count)
	binary.Read(f, binary.LittleEndian, ret)
	return ret
}

func LoadEdgeBasedEdges(filepath string) []graph.EdgeBasedEdge {
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

	ret := make([]graph.EdgeBasedEdge, count)
	binary.Read(f, binary.LittleEndian, ret)
	return ret
}
