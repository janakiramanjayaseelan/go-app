package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
    "fmt"
 "github.com/rollout/rox-go/v5/server"
	"github.com/gin-gonic/gin"
)


type Flags struct {
        EnableTutorial server.RoxFlag
}

var flags = &Flags{
        // Define the feature flags
        EnableTutorial: server.NewRoxFlag(false),
}

var rox *server.Rox



func main() {

     options := server.NewRoxOptions(server.RoxOptionsBuilder{})

        rox := server.NewRox()

        // Register the flags container with the CloudBees platform
        rox.RegisterWithEmptyNamespace(flags)

        // Setup the feature management environment key
        <-rox.Setup("781c4fb2-e537-47d9-6a10-7818cb92a644", options)

        // Boolean flag example
        fmt.Println("EnableTutorials value is " , flags.EnableTutorial.IsEnabled(nil))

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		sha := getCommitSha()
		color := getColor(sha)
		textColor := getTextColor(color)
		environment := getEnvironment()
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"sha":         sha,
			"color":       color,
			"textColor":   textColor,
			"environment": environment,
		})
	})

	router.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./static/favicon.ico")
	})

	router.Run()
}

func getCommitSha() string {
	content, err := ioutil.ReadFile("sha.txt")
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(content))
}

func getColor(sha string) string {
	h := sha256.New()
	h.Write([]byte(sha))
	hash := hex.EncodeToString(h.Sum(nil))
	return "#" + hash[:6]
}

func getTextColor(backgroundColor string) string {
	r, _ := strconv.ParseInt(backgroundColor[1:3], 16, 64)
	g, _ := strconv.ParseInt(backgroundColor[3:5], 16, 64)
	b, _ := strconv.ParseInt(backgroundColor[5:7], 16, 64)

	brightness := (r*299 + g*587 + b*114) / 1000

	if brightness > 155 {
		return "#000000"
	} else {
		return "#FFFFFF"
	}
}

func getEnvironment() string {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}
	return environment
}
