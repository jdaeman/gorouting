package extract_test

import (
	"extract"
	"fmt"
	"testing"
)

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
