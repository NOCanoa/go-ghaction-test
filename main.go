package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Githubstuff struct {
	Tag_name string `json:"tag_name"`
	Name     string `json:"name"`
	Body     string `json:"body"`
}

func main() {

	res, err := http.Get("https://api.github.com/repos/NOCanoa/go-ghaction-test/releases/latest")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic(fmt.Sprintf("response status error: %d %s", res.StatusCode, res.Status))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	/* fmt.Println(string(body)) */
	var github Githubstuff
	err = json.Unmarshal(body, &github)
	if err != nil {
		panic(err)
	}

	/* fmt.Print(github.Body) */

	webhookURL := "no-webhook"
	if len(os.Args) >= 2 {
		webhookURL = os.Args[1]
	}
	message := map[string]interface{}{
		"content":    "",
		"username":   "ZenDroid",
		"avatar_url": "https://cdn.jsdelivr.net/gh/zen-browser/www/public/logos/zen-alpha-black.png",
		"embeds": []map[string]interface{}{
			{
				"title":       "",
				"description": "# 🚀 New update " + github.Tag_name + "\r\n \r\n" + github.Body,
				"color":       16711680,
				"footer": map[string]interface{}{
					"text":     github.Tag_name,
					"icon_url": "https://cdn.jsdelivr.net/gh/zen-browser/www/public/logos/zen-alpha-black.png",
				},
			},
		},
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Fatalf("Failed to send message: %s", resp.Status)
	}

	log.Println("Message sent successfully!")
}
