module engine

go 1.17

require github.com/dhconnelly/rtreego v1.1.0

require graph v0.0.0

replace graph v0.0.0 => ../graph

require files v0.0.0

replace files v0.0.0 => ../files

require view v0.0.0

require github.com/jdaeman/go-shp v0.0.0-20190401125246-9fd306ae10a6 // indirect

replace view v0.0.0 => ../view
