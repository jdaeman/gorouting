package files

import (
	"encoding/binary"
	"errors"
	"modules/graph"
	"os"
)

// Common file format
// DataCount: 4 bytes
// DataList...

type Writer struct {
	savePath string
}

func NewWriter(savePath string) *Writer {
	ret := &Writer{
		savePath: savePath,
	}

	return ret
}

func (w Writer) createFile(ext string) (*os.File, error) {
	if w.savePath == "" {
		return nil, errors.New("empty path")
	}

	f, err := os.Create(ToDataPath(w.savePath, ext))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (w Writer) SaveGeoNodes(geoNodes []graph.ResultNode) error {
	f, err := w.createFile(GEONODE)
	if err != nil {
		return err
	}
	defer f.Close()

	err = binary.Write(f, binary.LittleEndian, int32(len(geoNodes)))
	if err != nil {
		return err
	}
	err = binary.Write(f, binary.LittleEndian, geoNodes)
	if err != nil {
		return err
	}
	return nil
}

func (w Writer) SaveEdges(edges []graph.InternalEdge) error {
	f, err := w.createFile(NBGEDGE)
	if err != nil {
		return err
	}
	defer f.Close()

	err = binary.Write(f, binary.LittleEndian, int32(len(edges)))
	if err != nil {
		return err
	}
	err = binary.Write(f, binary.LittleEndian, edges)
	if err != nil {
		return err
	}
	return nil

}

func (w Writer) SaveEdgeAnnotations(annotations []graph.EdgeAnnotation) error {
	f, err := w.createFile(ANNOTATION)
	if err != nil {
		return err
	}
	defer f.Close()

	err = binary.Write(f, binary.LittleEndian, int32(len(annotations)))
	if err != nil {
		return err
	}
	err = binary.Write(f, binary.LittleEndian, annotations)
	if err != nil {
		return err
	}
	return nil
}

func (w Writer) SaveEdgeGeometries(geometries []graph.Geometry) error {
	f, err := w.createFile(GEOMETRY)
	if err != nil {
		return err
	}
	defer f.Close()

	err = binary.Write(f, binary.LittleEndian, int32(len(geometries)))
	if err != nil {
		return err
	}

	for _, geometry := range geometries {
		err = binary.Write(f, binary.LittleEndian, int32(len(geometry.Nodes)))
		err = binary.Write(f, binary.LittleEndian, geometry.Nodes)
		err = binary.Write(f, binary.LittleEndian, geometry.Distances)
	}

	return err
}

func (w Writer) SaveTurnRestrictions(restrictions []graph.InternalRestriction) error {
	f, err := w.createFile(RESTRICTION)
	if err != nil {
		return err
	}
	defer f.Close()

	err = binary.Write(f, binary.LittleEndian, int32(len(restrictions)))
	if err != nil {
		return err
	}
	err = binary.Write(f, binary.LittleEndian, restrictions)
	if err != nil {
		return err
	}
	return nil
}

func (w Writer) SaveEdgeBasedNodes(nodes []graph.EdgeBasedNode) error {
	f, err := w.createFile(EBNODE)
	if err != nil {
		return err
	}
	defer f.Close()

	err = binary.Write(f, binary.LittleEndian, int32(len(nodes)))
	if err != nil {
		return err
	}
	err = binary.Write(f, binary.LittleEndian, nodes)
	if err != nil {
		return err
	}
	return nil
}

func (w Writer) SaveEdgeBasedEdges(edges []graph.EdgeBasedEdge) error {
	f, err := w.createFile(EBEDGE)
	if err != nil {
		return err
	}
	defer f.Close()

	err = binary.Write(f, binary.LittleEndian, int32(len(edges)))
	if err != nil {
		return err
	}
	err = binary.Write(f, binary.LittleEndian, edges)
	if err != nil {
		return err
	}
	return nil
}

func (w Writer) SaveEdgeBasedNodeSegments(segments []graph.EdgeBasedNodeSegment) error {
	f, err := w.createFile(EBSEG)
	if err != nil {
		return err
	}
	defer f.Close()

	err = binary.Write(f, binary.LittleEndian, int32(len(segments)))
	if err != nil {
		return err
	}
	err = binary.Write(f, binary.LittleEndian, segments)
	if err != nil {
		return err
	}
	return nil
}
