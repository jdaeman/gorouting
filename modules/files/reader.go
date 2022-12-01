package files

import (
	"encoding/binary"
	"errors"
	"modules/graph"
	"os"
)

type Reader struct {
	loadPath string
}

func NewReader(loadPath string) *Reader {
	ret := &Reader{
		loadPath: loadPath,
	}

	return ret
}

func (r Reader) openFile(ext string) (*os.File, error) {
	if r.loadPath == "" {
		return nil, errors.New("empty path")
	}

	f, err := os.Open(ToDataPath(r.loadPath, ext))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r Reader) LoadGeoNodes() []graph.ResultNode {
	f, err := r.openFile(GEONODE)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var count int32
	err = binary.Read(f, binary.LittleEndian, &count)

	if err != nil || count == 0 {
		if err != nil {
			panic(err)
		}
		panic("Geo node data count is zero")
	}

	ret := make([]graph.ResultNode, count)
	err = binary.Read(f, binary.LittleEndian, ret)
	if err != nil {
		panic(err)
	}
	return ret
}

func (r Reader) LoadEdges() []graph.InternalEdge {
	f, err := r.openFile(NBGEDGE)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var count int32
	err = binary.Read(f, binary.LittleEndian, &count)

	if err != nil || count == 0 {
		if err != nil {
			panic(err)
		}
		panic("edge data count is zero")
	}

	ret := make([]graph.InternalEdge, count)
	err = binary.Read(f, binary.LittleEndian, ret)
	if err != nil {
		panic(err)
	}
	return ret
}

func (r Reader) LoadEdgeAnnotations() []graph.EdgeAnnotation {
	f, err := r.openFile(ANNOTATION)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var count int32
	err = binary.Read(f, binary.LittleEndian, &count)

	if err != nil || count == 0 {
		if err != nil {
			panic(err)
		}
		panic("edge annotation data count is zero")
	}

	ret := make([]graph.EdgeAnnotation, count)
	err = binary.Read(f, binary.LittleEndian, ret)
	if err != nil {
		panic(err)
	}
	return ret
}

func (r Reader) LoadEdgeGeometries() []graph.Geometry {
	f, err := r.openFile(GEOMETRY)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var count, segCount int32
	err = binary.Read(f, binary.LittleEndian, &count)

	if err != nil || count == 0 {
		if err != nil {
			panic(err)
		}
		panic("edge geometry data count is zero")
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

func (r Reader) LoadRestrictions() []graph.InternalRestriction {
	f, err := r.openFile(RESTRICTION)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var count int32
	err = binary.Read(f, binary.LittleEndian, &count)

	if err != nil || count == 0 {
		if err != nil {
			panic(err)
		}
		panic("turn restriction data count is zero")
	}

	ret := make([]graph.InternalRestriction, count)
	err = binary.Read(f, binary.LittleEndian, ret)
	if err != nil {
		panic(err)
	}
	return ret
}

func (r Reader) LoadEdgeBasedNodes() []graph.EdgeBasedNode {
	f, err := r.openFile(EBNODE)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var count int32
	err = binary.Read(f, binary.LittleEndian, &count)

	if err != nil || count == 0 {
		if err != nil {
			panic(err)
		}
		panic("edge_based_node data count is zero")
	}

	ret := make([]graph.EdgeBasedNode, count)
	err = binary.Read(f, binary.LittleEndian, ret)
	if err != nil {
		panic(err)
	}
	return ret
}

func (r Reader) LoadEdgeBasedEdges() []graph.EdgeBasedEdge {
	f, err := r.openFile(EBEDGE)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var count int32
	err = binary.Read(f, binary.LittleEndian, &count)

	if err != nil || count == 0 {
		if err != nil {
			panic(err)
		}
		panic("edge_based_edge data count is zero")
	}

	ret := make([]graph.EdgeBasedEdge, count)
	err = binary.Read(f, binary.LittleEndian, ret)
	if err != nil {
		panic(err)
	}
	return ret
}

func (r Reader) LoadEdgeBasedNodeSegements() []graph.EdgeBasedNodeSegment {
	f, err := r.openFile(EBSEG)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var count int32
	err = binary.Read(f, binary.LittleEndian, &count)

	if err != nil || count == 0 {
		if err != nil {
			panic(err)
		}
		panic("edge_based_edge data count is zero")
	}

	ret := make([]graph.EdgeBasedNodeSegment, count)
	err = binary.Read(f, binary.LittleEndian, ret)
	if err != nil {
		panic(err)
	}
	return ret
}
