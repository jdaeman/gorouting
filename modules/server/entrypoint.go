package server

import (
	"log"
	"modules/engine"
	"modules/files"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type server_core struct {
	port int

	httpEngine  *gin.Engine
	routeEngine *engine.RoutingEngine
	repository  *files.DataRepository
}

func newServerCore(config ServerConfig) *server_core {
	httpEngine := gin.Default()

	log.Println("Start core data loading...")
	repository := files.NewDataRepository(config.DataPath)
	log.Println("Success!")
	log.Println("Node count", len(repository.GetNodes()))
	log.Println("Edge count", len(repository.GetEdges()))
	log.Println("Segment count", len(repository.GetSegments()))
	log.Println("Location count", len(repository.GetGeoNodes()))

	routingEngine := engine.NewRoutingEngineByData(
		repository.GetNodes(),
		repository.GetEdges(),
		repository.GetGeoNodes(),
		repository.GetSegments(),
	)

	ret := &server_core{
		port:        int(config.Port),
		httpEngine:  httpEngine,
		routeEngine: routingEngine,
		repository:  repository,
	}

	return ret
}

func (core *server_core) readyRouteService() {
	core.httpEngine.GET("/route/v1/driving/:coords", func(c *gin.Context) {
		reqCoords := c.Param("coords")
		coords, err := ParseCoordinates(reqCoords)

		reqStep := false
		steps, existParam := c.Get("steps")
		if existParam && steps.(bool) == true {
			reqStep = true
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"result":  "queryError",
				"message": err.Error(),
			})
			return
		}

		routeParam := engine.RouteParameter{Coords: coords, Steps: reqStep}
		var resp map[string]interface{}
		err = engine.Route(core.repository, core.routeEngine, routeParam, &resp)
		if err != nil || resp == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"result":  "fail",
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, resp)
		}
	})
}

func (core *server_core) startServer() bool {
	port := core.port

	core.readyRouteService()
	core.httpEngine.Run(":" + strconv.Itoa(port))
	return true
}

func Run(config ServerConfig) {
	server_core := newServerCore(config)
	server_core.startServer()
}
