package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

	r.GET("/", indexHandler)
	r.POST("/activateLevel", activateLevelHandler)

	r.Run(fmt.Sprintf(":%d", configuration.Server.Port))
}

func runClusterCommand(command string) string {
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
	}

	return string(output)
}

func indexHandler(c *gin.Context) {
	checkKubernetesClusterStatus()

	fmt.Printf("--- l:%+v\n", levels[0])

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Levels": levels,
		"Config": configuration,
	})
}

func activateLevelHandler(c *gin.Context) {
	levelId := c.PostForm("levelId")

	l, err := getLevelById(levelId)
	if err != nil {
		strMessage := "No level was found at: " + levelId
		errors.New(strMessage)
	}

	fmt.Println(&l)

	l.Active = true //TODO: set value
	var output = applyLevelResources(l.ResourcesPath)

	c.HTML(http.StatusOK, "applyOutput.tmpl", gin.H{
		"Output": output,
	})
}

func getLevelById(Id string) (*level.Level, error) {
	for _, level := range levels {
		if level.Id.String() == Id {
			return &level, nil
		}
	}

	return nil, errors.New("Level was not found")
}

func checkKubernetesClusterStatus() {
	output := runClusterCommand("kubectl cluster-info")
	if strings.Contains(output, "is running at") {
		configuration.Server.ClusterUp = true
	} else {
		configuration.Server.ClusterUp = false
	}
}

func applyLevelResources(path string) string {
	output := runClusterCommand("kubectl apply -f " + path + " --recursive")

	return output
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
		l.Id = uuid.New()
		l.ResourcesPath = filepath.Dir(matches[p]) + "/resources/"

		err = yaml.Unmarshal(f, &l)
		levels = append(levels, l)

		if err != nil {
			log.Fatalf("error: %v", err)
		}

		fmt.Printf("--- l:%+v\n", l)
	}
}
