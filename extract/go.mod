module extract

go 1.17

require github.com/paulmach/osm v0.7.0

require github.com/paulmach/orb v0.1.3 // indirect

require util v0.0.0

replace util v0.0.0 => ../util

require graph v0.0.0

replace graph v0.0.0 => ../graph

require files v0.0.0

replace files v0.0.0 => ../files

//require github.com/dhconnelly/rtreego v1.1.0