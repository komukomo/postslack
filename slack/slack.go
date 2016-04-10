package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type SlackMessage struct {
	Text        string       `json:"text"`
	Username    string       `json:"username"`
	Channel     string       `json:"channel"`
	Icon        string       `json:"icon_emoji"`
	Attachments []Attachment `json:"attachments"`
}

type Config struct {
	Url     string `json:"url"`
	Channel string `json:"channel"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
}

type Attachment struct {
	Fallabck   string  `json:"fallback`
	Color      string  `json:"color"`
	Pretext    string  `json:"pretext"`
	AuthorName string  `json:"author_name"`
	Title      string  `json:"title"`
	Text       string  `json:"text"`
	Fields     []Field `json:"fields"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func PostSlack(incomingURL string, msg SlackMessage) {
	params, _ := json.Marshal(msg)
	values := url.Values{"payload": {string(params)}}
	Post(incomingURL, values)
}

func Post(url string, values url.Values) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		fmt.Print(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()
}

func main() {
	incomingURL := os.Getenv("SLACK_URL")
	PostSlack(incomingURL, SlackMessage{
		"message",
		"botname",
		"#random",
		":ghost:",
		nil,
	})
}
