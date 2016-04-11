package main

import (
	"./slack"
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func main() {
	config := loadDefaultConfig("./config.json")

	args := make(map[string]*string)
	args["channel"] = flag.String("c", config.Channel, "cannel name")
	args["botname"] = flag.String("n", config.Name, "bot name")
	args["icon"] = flag.String("i", config.Icon, "bot icon. emoji or URL ")
	args["incomingURL"] = flag.String("url", config.Url, "incomingURL")
	args["attachmentsFile"] = flag.String("a", "", "attachment filepath")
	args["param"] = flag.String("p", "", "parameters")
	noStdin := flag.Bool("e", false, "no stdin (for attachments post)")
	flag.Parse()

	output := ""
	if !*noStdin {
		output = getStdin()
	}

	if *args["attachmentsFile"] != "" {
		parameters := str2map(*args["param"], output)
		postAttachments(*args["incomingURL"], *args["attachmentsFile"], parameters)
	} else {
		simplePost(args, output)
	}
}

func postAttachments(incomingURL string, attachmentsFile string, parameters map[string]string) {
	var doc bytes.Buffer
	slackMessage := slack.SlackMessage{}
	tpl := template.Must(template.ParseFiles(attachmentsFile))
	tpl.Execute(&doc, parameters)
	json.Unmarshal(doc.Bytes(), &slackMessage)
	slack.PostSlack(incomingURL, slackMessage)
}

func simplePost(args map[string]*string, text string) {
	slack.PostSlack(*args["incomingURL"], slack.SlackMessage{
		text,
		*args["botname"],
		*args["channel"],
		*args["icon"],
		nil,
	})
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func str2map(param string, str string) map[string]string {
	paramlist := strings.Split(param, "&")
	result := make(map[string]string)

	for _, val := range paramlist {
		a := strings.Split(val, "=")
		if a[1] == "__stdin" {
			result[a[0]] = strings.Replace(str, "\n", "\\n", -1)
		} else {
			result[a[0]] = a[1]
		}
	}
	return result
}

func getStdin() (stdin string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		stdin += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading stdin:", err)
	}
	return stdin
}

func loadDefaultConfig(configFilePath string) slack.Config {
	config := slack.Config{
		os.Getenv("SLACK_URL"),
		os.Getenv("SLACK_CHANNEL"),
		os.Getenv("SLACK_ICON"),
		os.Getenv("SLACK_NAME"),
	}

	if exists(configFilePath) {
		file, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(file, &config)
	}
	return config
}
