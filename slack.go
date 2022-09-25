package main

import (
  "os"
  "log"
  "gopkg.in/yaml.v2"
  "fmt"
  slk "github.com/slack-go/slack"
)

func post_slack(user string, text string) {
    f, err := os.Open("speak.cfg")
    decoder := yaml.NewDecoder(f)
    err = decoder.Decode(&cfg)
    if (err != nil) {
      log.Fatal("Config Decode error: ",err)
    }

    api := slk.New(cfg.BotToken)
   	c, rw, _, err := api.JoinConversation(cfg.NotifyChannel)
	if err != nil {
		fmt.Printf("Join Error %s %s\n", rw,err)
		return
	}
  fmt.Printf("Got channel %v\n",c)
  _, _, err = api.PostMessage(cfg.NotifyChannel,slk.MsgOptionText("This is a test",false))
}
