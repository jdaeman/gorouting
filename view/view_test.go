package view_test

import (
	"testing"

	"github.com/jdaeman/go-shp"
)

func TestPointShp(t *testing.T) {
	w, err := shp.Create("1_point", shp.POINT)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	point := &shp.Point{X: 126.9996771, Y: 37.5589237}
	w.Write(point)
}

func TestLineShp(t *testing.T) {
	w, err := shp.Create("1_polyline", shp.POLYLINE)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	polyline := shp.NewPolyLine([][]shp.Point{
		{{X: 127.0011579, Y: 37.5856008}, {X: 127.0029363, Y: 37.5872848}},
		{{X: 127.0029363, Y: 37.5872848}, {X: 127.0056584, Y: 37.5884409}},
		{{X: 127.0056584, Y: 37.5884409}, {X: 127.0084712, Y: 37.5895794}},
	})

	w.Write(polyline)
}

// func TestView(t *testing.T) {

// 	line := make([]shp.Point, 0, 2)

// 	line = append(line, shp.Point{X: 126.9996771, Y: 37.5589237})
// 	line = append(line, shp.Point{X: 127.0028535, Y: 37.5852564})

// 	polyline := shp.NewPolyLine([][]shp.Point{
// 		line,
// 	})

// 	// fields := []shp.Field{
// 	// 	// String attribute field with length 25
// 	// 	shp.StringField("NAME", 25),
// 	// }

// 	w, err := shp.Create("abc", shp.POLYLINE)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// .dbf 로 수정해야 함.
// 	defer w.Close()

// 	w.Write(polyline)
// 	//w.SetFields(fields)
// 	//r := w.Write(&shp.Point{X: 126.9996771, Y: 37.5589237})
// 	//w.WriteAttribute(0, 0, "ABCC")
// }
