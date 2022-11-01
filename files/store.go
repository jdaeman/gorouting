package files

import (
	"encoding/binary"
	"graph"
	"os"
)

// Common file format
// DataCount: 8bytes
// DataList...

// Store node locations
func StoreGeoNodes(filepath string, geoNodes []graph.ExternalNode) error {
	f, err := os.Create(filepath)
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

// Store uncompressed edges
func StoreUncompEdges(filepath string, edges []graph.InternalEdge) error {
	f, err := os.Create(filepath)
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

// Store edge annotation.(from osm way)
func StoreEdgeAnnotations(filepath string, annotations []graph.EdgeAnnotation) error {
	f, err := os.Create(filepath)
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
