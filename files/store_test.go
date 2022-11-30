package files_test

import (
	"files"
	"fmt"
	"graph"
	"testing"
)

func TestGeoNodes(t *testing.T) {

	data := make([]graph.ResultNode, 2)

	data[0].Id = 88
	data[0].X = 120.123
	data[0].Y = 36.123

	data[1].Id = 99
	data[1].X = 120.321
	data[1].Y = 36.321

	files.StoreGeoNodes("test.geo", data)

	load := files.LoadGeoNodes("test.geo")

	if len(load) != 2 {
		t.FailNow()
	}

	fmt.Println(load)
}

func TestOutputPath(t *testing.T) {
	path := "data\\osm\\map.osm"
	newPath := files.ToDataPath(path, ".geos")

	fmt.Println("newPath", newPath)
	if newPath != "data\\osm\\map.geos" {
		t.Fail()
	}
}
