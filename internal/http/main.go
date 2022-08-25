package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	config "learn-k8s/internal/config"
	level "learn-k8s/internal/level"
)

var (
	levels        []level.Level
	configuration config.Config
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("internal/http/templates/*")

	LoadConfig()
	LoadLevels()

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Title":  "Main website",
			"Levels": levels,
		})
	})

	cmd := exec.Command("kubectl", "get", "pods")
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
	}

	fmt.Print(string(output))

	ginPort := configuration.Server.Port
	r.Run(fmt.Sprintf(":%d", ginPort))
}

func LoadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./internal/http")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
}

func LoadLevels() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	dir := fmt.Sprintf("%s/levels/*/level.yaml", path)
	matches, err := filepath.Glob(dir)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for p := range matches {
		fmt.Println(matches[p])

		f, err := os.ReadFile(matches[p])
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		var l level.Level

		err = yaml.Unmarshal(f, &l)
		levels = append(levels, l)

		if err != nil {
			log.Fatalf("error: %v", err)
		}

		fmt.Printf("--- l:%+v\n", l)
	}
}
