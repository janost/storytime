//go:generate pkger

package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/janost/storytime/model"
	"github.com/kelseyhightower/envconfig"
)

type Specification struct {
	OllamaHost  string
	OllamaPort  string
	OllamaModel string
	DbHost      string
	DbName      string
	DbUser      string
	DbPass      string
}

type OllamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaGenerateResponse struct {
	Model      string
	Created_at string
	Response   string
}

type StoryRequest struct {
	Cast     string `json:"cast"`
	Location string `json:"location"`
	Plot     string `json:"plot"`
}

var ollamaUrl string
var s Specification

//go:embed views/*
var viewsfs embed.FS

func main() {
	err := envconfig.Process("storytime", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Using ollama at %s:%s", s.OllamaHost, s.OllamaPort)
	ollamaUrl = fmt.Sprintf("http://%s:%s/api/generate", s.OllamaHost, s.OllamaPort)
	model.Setup(s.DbHost, s.DbName, s.DbUser, s.DbPass)
	// create new fiber instance  and use across whole app
	engine := html.NewFileSystem(http.FS(viewsfs), ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// middleware to allow all clients to communicate using http and allow cors
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		stories, _ := model.GetAllStories()
		// Render index - start with views directory
		return c.Render("views/index", fiber.Map{
			"Stories": stories,
		})
	})

	app.Post("/create", handleCreate)

	log.Fatal(app.Listen(":4000"))
}

func handleCreate(c *fiber.Ctx) error {
	storyRequest := new(StoryRequest)

	if err := c.BodyParser(storyRequest); err != nil {
		fmt.Println("error = ", err)
		return c.SendStatus(200)
	}

	storyPrompt := "Write a story suitable for a small child. The story should be compelling, cute and proper for small children, and be as long as you can write.\n"
	if len(storyRequest.Cast) > 0 {
		storyPrompt += fmt.Sprintf("The story's main characters are: %s.\n", storyRequest.Cast)
	}
	if len(storyRequest.Location) > 0 {
		storyPrompt += fmt.Sprintf("The story's location is: %s.\n", storyRequest.Location)
	}
	if len(storyRequest.Plot) > 0 {
		storyPrompt += fmt.Sprintf("The story's plot is: %s.", storyRequest.Plot)
	}

	ollamaRequest, _ := json.Marshal(OllamaGenerateRequest{
		Model:  s.OllamaModel,
		Prompt: storyPrompt,
		Stream: false,
	})

	fmt.Println("ollamaRequest = ", string(ollamaRequest))

	client := resty.New()
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(string(ollamaRequest)).
		Post(ollamaUrl)

	doc := OllamaGenerateResponse{}
	json.Unmarshal(resp.Body(), &doc)

	model.CreateStory(storyRequest.Cast, storyRequest.Location, storyRequest.Plot, doc.Response, s.OllamaModel)

	c.Response().Header.Add("HX-Redirect", "/")

	return c.Redirect("/")
}
