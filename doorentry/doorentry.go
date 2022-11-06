package main

import (
  "log"
  "fmt"
  "net/http"
  "net/url"
  "database/sql"
  "encoding/base64"
  "os"
  "flag"
  "strconv"
   "gopkg.in/yaml.v2"
  _ "modernc.org/sqlite"
)

type DoorConfig struct {
   SecretKey string `yaml:"SecretKey"`
   AccessKey string `yaml:"AccessKey"`
   DBFile string `yaml:"DBFile"`
   AudioDir string `yaml:"AudioDir"`
   BottomSpeaks []string `yaml:"BottomSpeaks"`
}

var cfg DoorConfig

func main() {
  configFile := flag.String("config","doorentry.cfg","Path to config file")
  flag.Parse()

  f, err := os.Open(*configFile)
  decoder := yaml.NewDecoder(f)
  err = decoder.Decode(&cfg)
  if (err != nil) {
    log.Fatal("Config Decode error: ",err)
  }

  memberId,err := strconv.Atoi(flag.Arg(0))
  if (err != nil) { log.Fatal("Invalid Member ID sepcified: ",err) }
  db, err := sql.Open("sqlite", cfg.DBFile)
  if (err != nil) { log.Fatal("Cannot open DB File",cfg.DBFile,err) }
  defer func() { db.Close() }()
  stmt, err := db.Prepare("SELECT member, firstname, lastname, nickname FROM members WHERE id = ?")
  if (err != nil) { log.Fatal("Prepare failed ",err) }
  var name,member,first,last,nick string
  err = stmt.QueryRow(memberId).Scan(&member,&first, &last, &nick)
  if (err != nil) { log.Fatal("No member found ",name,err) }
  if (nick =="") {
    name = first+" "+last
  } else {
    name = nick+" "+last
  }
  fmt.Println(name)

  pcmstr := ""

  filename := cfg.AudioDir+"/"+member+".pcm"
  _, err = os.Stat(filename)
  if (err  != nil) {
      log.Println("Member audio file non existant ",name,err) 
  } else 
  {
    pcmdata, err := os.ReadFile(filename)
    if (err != nil) { 
      log.Println("Error reading member PCM data ",name,err) 
    } else {
      pcmstr = base64.URLEncoding.EncodeToString(pcmdata)
    }
  }
  for _,bs := range cfg.BottomSpeaks {
        fmt.Fprintf(os.Stderr, "Dispatch to Bottom: \"%s\"\n",bs)
        response, err := http.PostForm(bs, url.Values{
        "quickText": { "Welcome "+name},
        "audio": { pcmstr }})
        if (err != nil) {
          fmt.Fprintf(os.Stderr, "Bottom response from %s is %v %s\n",bs,response,err)
        }
      }

}

