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
func StoreGeoNodes(filepath string, geoNodes []graph.ResultNode) error {
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
func StoreEdges(filepath string, edges []graph.NodeBasedEdge) error {
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

func StoreEdgeGeometries(filepath string, geometries []graph.Geometry) error {
	f, err := os.Create(filepath)
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

func StoreTurnRestrictions(filepath string, restrictions []graph.InternalRestriction) error {
	f, err := os.Create(filepath)
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
