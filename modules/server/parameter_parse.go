package server

import (
	"errors"
	"strconv"
	"strings"
)

func ParseCoordinates(coordinates string) ([][2]float64, error) {
	coordList := strings.Split(coordinates, ";")

	ret := make([][2]float64, 0, len(coordList))

	for _, coords := range coordList {
		xy := strings.Split(coords, ",")

		if len(xy) != 2 {
			return nil, errors.New("coordinate error " + coords)
		}

		strX, strY := xy[0], xy[1]

		x, err1 := strconv.ParseFloat(strX, 64)
		y, err2 := strconv.ParseFloat(strY, 64)

		if err1 != nil || err2 != nil {
			return nil, errors.New("Parse error " + strX + " " + strY)
		}

		ret = append(ret, [2]float64{x, y})
	}

	return ret, nil
}
