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
   	_, _, _, err = api.JoinConversation(cfg.NotifyChannel)
  if (err != nil) {
    fmt.Fprintf(os.Stderr,"Slack Join Channel error %s\n",err)
  }
  //fmt.Printf("Got channel %v\n",c)
  _, _, err = api.PostMessage(cfg.NotifyChannel,slk.MsgOptionText(fmt.Sprintf("%s posted \"%s\"",user,text),false))
  if (err != nil) {
    fmt.Fprintf(os.Stderr,"Slack post message error %s\n",err)
  }
}
