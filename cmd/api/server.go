package api

import (
	"Test_task_Halolab/cmd/datagen"
	"Test_task_Halolab/cmd/service"
	_ "Test_task_Halolab/docs"
	"Test_task_Halolab/pkg/database"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	"runtime"
	"time"
)

type Server struct {
	Api   *API
	Redis *redis.Client
}

func NewServer() Server {
	return Server{
		Api: NewAPI(),
	}
}

func (s *Server) StartServer(ctx context.Context) error {
	// init application layers
	err := s.InitLayers()
	if err != nil {
		return err
	}

	// Start data generator
	withCancel, cancelFunc := context.WithCancel(ctx)
	dataGen := datagen.NewDataGen(withCancel)
	dataGen.StartSensorGroups()
	defer dataGen.WaitGroup.Wait()
	defer cancelFunc()

	// Save groups to db
	err = s.Api.Service.SaveSensorGroup(dataGen.Groups)
	if err != nil {
		return err
	}

	defer func() {
		if err := s.Api.Service.Redis.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	// Start data aggregator goroutines
	errChan := s.Api.Service.StartDataAggregator(ctx, dataGen.Pool)
	go func() {
		for {
			select {
			case <-ctx.Done():
				cancelFunc()
				return
			case err := <-errChan:
				fmt.Println(err)
			default:
				continue
			}
		}
	}()

	// Start http server
	srv := http.Server{
		Addr:    ":8080",
		Handler: s.RegisterRoutes(),
	}

	ch := make(chan error, 1)
	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
	}()

	// wait for error or done ctx
	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		close(dataGen.Pool)
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		return srv.Shutdown(timeout)
	}
}

func (s *Server) InitLayers() error {
	// init DB layer
	db := database.NewDatabase()
	err := db.InitDatabase()
	if err != nil {
		return err
	}

	// init service layer
	serviceLayer := service.NewService()
	err = serviceLayer.InitService(db)
	if err != nil {
		return err
	}

	// init api layer
	s.Api.InitAPI(serviceLayer)

	return nil
}

func (s *Server) RegisterRoutes() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		runtime.Gosched()
		c.JSON(http.StatusOK, "pong")
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger/doc.json")))

	groupRoutes := router.Group("/group/:groupName")
	{
		groupRoutes.GET("/temperature/average", s.Api.GetGroupAverageTemperature)
		groupRoutes.GET("/transparency/average", s.Api.GetGroupAverageTransparency)
		groupRoutes.GET("/species", s.Api.GetGroupSpecies)
		groupRoutes.GET("/species/top/:N", s.Api.GetTopNSpecies)
	}

	regionRoutes := router.Group("/region")
	{
		regionRoutes.GET("/temperature/min", s.Api.GetRegionTemperatureMin)
		regionRoutes.GET("temperature/max", s.Api.GetRegionTemperatureMax)
	}

	sensorRoutes := router.Group("/sensor")
	{
		sensorRoutes.GET("/:codeName/temperature/average", s.Api.GetSensorAverageTemperature)
	}

	return router
}
