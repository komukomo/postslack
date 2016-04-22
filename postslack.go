package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/komukomo/postslack/slack"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type PostSlack struct {
}

func (ps PostSlack) Run(osArgs []string) {
	configFilePath := os.Getenv("HOME") + "/.postslackrc"
	attachmentsFilePath := os.Getenv("HOME") + "/.at-postslackrc"
	config := loadDefaultConfig(configFilePath)

	flags := flag.NewFlagSet("postslack", flag.ContinueOnError)

	args := make(map[string]*string)
	args["channel"] = flags.String("ch", config.Channel, "channel name")
	args["botname"] = flags.String("name", config.Name, "bot name")
	args["icon"] = flags.String("icon", config.Icon, "bot icon. emoji or URL ")
	args["incomingURL"] = flags.String("url", config.Url, "incomingURL")
	args["attachmentsFile"] = flags.String("att", attachmentsFilePath, "attachment filepath")
	args["param"] = flags.String("param", "", "parameters")
	noStdin := flags.Bool("empty", false, "no stdin (for attachments post)")
	noAttachments := flags.Bool("no-attachments", false, "no attachments")
	flags.Parse(osArgs)

	if *args["incomingURL"] == "" {
		panic("no value for incoming-webhook URL")
	}

	output := ""
	if !*noStdin {
		output = getStdin()
	}

	if *noAttachments {
		simplePost(args, output)
	} else {
		if exists(*args["attachmentsFile"]) {
			parameters := str2map(*args["param"], output)
			postAttachments(*args["incomingURL"], *args["attachmentsFile"], parameters)
		} else {
			panic("no attachmentsFile")
		}
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

	result["Stdin"] = strings.Replace(str, "\n", "\\n", -1)
	for _, val := range paramlist {
		a := strings.Split(val, "=")
		if len(a) == 2 {
			if a[1] == "__stdin" {
				result[a[0]] = strings.Replace(str, "\n", "\\n", -1)
			} else {
				result[a[0]] = a[1]
			}
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
	config := loadEnvConfig()

	if exists(configFilePath) {
		file, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(file, &config)
	}
	return config
}

func loadEnvConfig() slack.Config{
	return slack.Config{
		os.Getenv("SLACK_URL"),
		os.Getenv("SLACK_CHANNEL"),
		os.Getenv("SLACK_ICON"),
		os.Getenv("SLACK_NAME"),
	}
}
