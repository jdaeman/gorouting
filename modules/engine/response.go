package engine

type respWayPoint struct {
	Name     string     `json:"name"`
	Location [2]float64 `json:"location"`
}

type respManeuver struct {
	BearingAfter  int        `json:"bearing_after"`
	BearingBefore int        `json:"bearing_before"`
	Type          string     `json:"type"`     // depart or arrive
	Location      [2]float64 `json:"location"` // [0]: x, [1]: y
}

type respIntersection struct {
	Entry    [1]bool    `json:"entry"`    // [true]
	Bearings [1]int     `json:"bearings"` // [10]
	Location [2]float64 `json:"location"` // {x, y}
}

type respStep struct {
	Distance      int                 `json:"distance"`
	Duration      int                 `json:"duration"`
	Weight        int                 `json:"weight"`
	Geometry      string              `json:"geometry"`
	Maneuver      respManeuver        `json:"maneuver"`
	Intersections [1]respIntersection `json:"intersections"`

	// DrivingSide string `json:"driving_side"`
	// Mode        string `json:"mode"`
	// Name        string `json:"name"`
}

type respLeg struct {
	Distance int        `json:"distance"`
	Duration int        `json:"duraiton"`
	Weight   int        `json:"weight"`
	Summary  string     `json:"summary"`
	Steps    []respStep `json:"steps"`
}

type respRoute struct {
	Distance   int       `json:"distance"`
	Duration   int       `json:"duration"`
	Weight     int       `json:"weight"`
	WeightName string    `json:"weight_name"`
	Legs       []respLeg `json:"legs"`
}
