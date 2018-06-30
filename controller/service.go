package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

func Run(s ServiceController) {
	g := gin.New()

	g.POST("/v1/register", func(c *gin.Context) {
		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, fmt.Errorf("read request body: %v", err))
			return
		}
		defer c.Request.Body.Close()

		var in RegisterInput
		if err := json.Unmarshal(b, &in); err != nil {
			c.JSON(500, fmt.Errorf("unmarshal request body: %v", err))
			return
		}

		out := s.Register(&in)
		c.JSON(out.Status, out)
	})

	g.GET("/v1/service", func(c *gin.Context) {
		out := s.Service()
		c.JSON(out.Status, out)
	})

	g.GET("/v1/service/:service_id", func(c *gin.Context) {
		in := &CatalogInput{
			ServiceID: c.Param("service_id"),
		}
		catalog := s.Catalog(in)
		c.JSON(catalog.Status, catalog)
	})

	g.GET("/v1/instance", func(c *gin.Context) {
		out := s.Instance()
		c.JSON(out.Status, out)
	})

	g.GET("/v1/instance/:instance_id", func(c *gin.Context) {
		in := &DescribeInput{
			InstanceID: c.Param("instance_id"),
		}
		out := s.Describe(in)
		c.JSON(out.Status, out)
	})

	g.POST("/v1/instance", func(c *gin.Context) {
		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, fmt.Errorf("read request body: %v", err))
			return
		}
		defer c.Request.Body.Close()

		var in CreateInput
		if err := json.Unmarshal(b, &in); err != nil {
			c.JSON(500, fmt.Errorf("unmarshal request body: %v", err))
			return
		}

		out := s.Create(&in)
		c.JSON(out.Status, out)
	})

	log.Printf("config=%v\n", s.Config())
	if err := g.Run(s.Config().Port); err != nil {
		log.Fatalf("run broker: %v", err)
	}
}
