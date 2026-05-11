package main

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Apps []App `yaml:"apps"`
	Version string `yaml:"version"`
	Date string `yaml:"date"`
}
func (c Config) String() string {
	return fmt.Sprintf(
		"Version: %v\tDate: %v",
		c.Version,
		c.Date,
	)
}

type App struct {
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	Target string `yaml:"target"`
	Version string `yaml:"version"`
	Description string `yaml:"description"`
}

func (a App) String() string {
	spaces := 15
	return fmt.Sprintf(
		"%-*s %v\n%-*s %v\n%-*s %v\n%-*s %v",
		spaces, "Name:", a.Name,
		spaces, "Description:", a.Description,
		spaces, "Path:", a.Path,
		spaces, "Version:", a.Version,
	)
}

func proxy(target string) gin.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	return func(c *gin.Context) {
		c.Request.URL.Path = c.Param("proxyPath")
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func readConfig(config *Config) bool{
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return false
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return false
	}
	return true
}

func main(){
	var config Config
	success := readConfig(&config)
	if !success && len(config.Apps) == 0 {
		log.Fatal("Failed to read configuration")
	}

	gw := gin.Default()
	log.Println("Gateway: ", config.String())
	log.Println("Gateway started with the following configuration:")
	for _, app := range config.Apps {
		gw.Any(app.Path+"/*proxyPath", proxy(app.Target))
		log.Println(app.String())
	}
	log.Println("Gateway listening at port :80")
	gw.Run(":80")
}
