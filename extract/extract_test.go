package extract_test

import (
	"extract"
	"fmt"
	"sort"
	"testing"
)

func TestExtractor(t *testing.T) {
	config := "data/map.osm"
	datas, err := extract.ReadOSM(config)

	if err != nil {
		t.Fail()
	}

	extractor := extract.NewExtractor(datas[0], datas[1], datas[2])
	extractor.ProcessOSMNodes()
	extractor.ProcessOSMWays()

	fmt.Println("OSM Node count", len(extractor.AllNodes))
	fmt.Println("OSM Edge count", len(extractor.AllEdges))

	extractor.ProcessNodes()
	extractor.ProcessEdges()

	fmt.Println("node count", len(extractor.UsedNodes))
	fmt.Println("edge count", len(extractor.UsedEdges))

	extractor.PrepareData()

	fmt.Println("Unique loc count", len(extractor.AllNodes))
	fmt.Println("Unique node count", len(extractor.UsedNodes))

}

func TestReadOSM(t *testing.T) {
	config := "data/map.osm"
	datas, err := extract.ReadOSM(config)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if len(datas[0]) == 0 {
		fmt.Println("Empty node")
	} else {
		fmt.Println("Nodes", len(datas[0]))
	}

	if len(datas[1]) == 0 {
		fmt.Println("Empty way")
	} else {
		fmt.Println("Ways", len(datas[1]))
	}

	if len(datas[2]) == 0 {
		fmt.Println("Empty restriction")
	} else {
		fmt.Println("Restrictions", len(datas[2]))
	}
}

func TestParseOSMData(t *testing.T) {
	config := extract.Config{"data/map.osm"}

	extract.ParseOSMData(config)
}

func TestUnique(t *testing.T) {
	dataList := []int{1, 1, 1, 2, 2, 3, 4, 5, 5}

	newIndex, first := 0, dataList[0]
	tmp := dataList
	for len(tmp) > 0 {
		newIndex += 1

		idx := sort.Search(len(tmp), func(idx int) bool {
			return tmp[idx] > first
		})

		if idx >= len(tmp) {
			tmp = nil
		} else {
			dataList[newIndex] = tmp[idx]
			first = tmp[idx]
			tmp = tmp[idx:]
		}
	}

	if newIndex != 5 {
		t.Fail()
	}

	for i := 0; i < newIndex; i++ {
		if dataList[i] != i+1 {
			t.Fail()
		}
	}

	fmt.Println(dataList[:newIndex])
}

func TestUnique1(t *testing.T) {
	dataList := []int{1, 2, 3, 4, 5}

	newIndex, first := 0, dataList[0]
	tmp := dataList
	for len(tmp) > 0 {
		newIndex += 1

		idx := sort.Search(len(tmp), func(idx int) bool {
			return tmp[idx] > first
		})

		if idx >= len(tmp) {
			tmp = nil
		} else {
			dataList[newIndex] = tmp[idx]
			first = tmp[idx]
			tmp = tmp[idx:]
		}
	}

	if newIndex != 5 {
		t.Fail()
	}

	for i := 0; i < newIndex; i++ {
		if dataList[i] != i+1 {
			t.Fail()
		}
	}
	fmt.Println(dataList[:newIndex])
}
