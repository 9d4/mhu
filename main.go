package main

import (
	"bytes"
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
)

var (
	ADDRESS      = "0.0.0.0:4500"
	MIKROTIK_URI = "https://admin:@192.168.104.12/rest/"

	MIKROTIK_URL *url.URL
)

var defaultTransport = http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

func init() {
	if addr := os.Getenv("ADDRESS"); addr != "" {
		ADDRESS = addr
	}

	if rest := os.Getenv("MIKROTIK_URI"); rest != "" {
		MIKROTIK_URI = rest
	}

	url, err := url.Parse(MIKROTIK_URI)
	if err != nil {
		log.Fatal(err)
	}

	MIKROTIK_URL = url
}

func main() {
	app := fiber.New()

	api := app.Group("/api")
	api.Get("/userprofiles", proxyGet("/ip/hotspot/user/profile"))
	api.Get("/servers", proxyGet("/ip/hotspot"))
	api.Post("/upload", handleSubmitUsers)

	app.Static("/", "views/dist/", fiber.Static{
		Browse: false,
		Index:  "index.html",
	})

	app.Listen(ADDRESS)
}

func getRestUri(path string) string {
	return MIKROTIK_URL.JoinPath(path).String()
}

func proxyGet(target string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cl := http.Client{Transport: &defaultTransport}
		res, err := cl.Get(getRestUri(target))
		if err != nil {
			log.Println(err)
			return err
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		c.Response().Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		c.Write(b)
		return nil
	}
}

func handleSubmitUsers(c *fiber.Ctx) error {
	server, profile := c.FormValue("server"), c.FormValue("profile")

	if server == "" || profile == "" {
		return fiber.ErrBadRequest
	}

	f, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	file, err := f.Open()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	for _, v := range data {
		if len(v) < 2 {
			continue
		}

		if v[0] == "" || v[1] == "" {
			continue
		}

		err = pushUser(server, profile, v[0], v[1])
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func pushUser(server, profile, name, password string) error {
	cl := http.Client{Transport: &defaultTransport}

	user := struct {
		Server   string `json:"server"`
		Profile  string `json:"profile"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}{server, profile, name, password}

	userjson := bytes.Buffer{}

	if err := json.NewEncoder(&userjson).Encode(user); err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	res, err := cl.Post(getRestUri("/ip/hotspot/user/add"), fiber.MIMEApplicationJSON, &userjson)
	if err != nil {
		log.Println(err)
		return err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	log.Println(string(b))

	return nil
}
