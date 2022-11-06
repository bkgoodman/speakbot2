package main

/*
Writing packet: b'\x00\x00\x00\x00\x00\x01Z00\x02E$\x04'
Writing packet: b'\x00\x00\x00\x00\x00\x01Z00\x02E$1BL000E0000AAU0040FFFF1AU0064FEFE2AU0064FEFE3AU0064FEFE4AU0064FEFE5AU0064FEFE\x04'
Writing packet: b'\x00\x00\x00\x00\x00\x01Z00\x02E.TUA\x04'
Writing packet: b'\x00\x00\x00\x00\x00\x01Z00\x02G1\x04'
Writing packet: b'\x00\x00\x00\x00\x00\x01Z00\x02AA\x1b bTesting 1-2-3\x04'
*/

/*

speak.cfg
echo 'U2FsdGVkX19Zh5RWjqNMl5n7AkOSCCSEbdqymT2IcP6+43QNTNP9cM94ltRX2wYuCvbGGRh7WtnDN1ujedMtuXOPMJ4T1TZ61Gb64rJD2jckf7HKkzGveEs8JVZxtwbt0dUE17LWxipWxl78zzfTTwuYXcj4NJEFk7DExVJ28wzLT6FtwItR/vBfCISsHddHWZpI6E94zdVqbH2mNlrpi5f/s8tbTQKBRgPbXVmqhJBCXXvDv5I7ZsUEgC89c0di4t/hsP7q0r6Cx7gop3vn81HzRo2+XhZHP4tqHyWBw4w=' | openssl enc -aes-256-cbc -md sha512 -pbkdf2 -A -a -d


 aplay -D sysdefault:CARD=PCH -c 2 -f S16_LE  -r 22050  data.pcm

Alsa device: sysdefault:CARD=PCH

See devices w/ aplay -L

*/

import (
   "fmt"
   "golang.org/x/exp/slices"
   "net/http"
   "net/http/cgi"
   "net/url"
   "bytes"
   "context"
   "time"
   "os"
   "os/exec"
   "encoding/base64"
   "gopkg.in/yaml.v2"
   //"syscall"

  "log"
  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/aws/aws-sdk-go-v2/config"
  "github.com/aws/aws-sdk-go-v2/service/polly"
  pollytype "github.com/aws/aws-sdk-go-v2/service/polly/types" 
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// Use for channel to pass speak request
type SpeakRequest struct {
  Text string 
  SlashCommand string
}

type BottomSpeak struct {
   URL string `yaml:"URL"`
   Commands []string `yaml:"Commands"`
}
type SpeakConfig struct {
   SecretKey string `yaml:"SecretKey"`
   AccessKey string `yaml:"AccessKey"`
   BotToken string `yaml:"BotToken"`
   Token string `yaml:"Token"`
   Port int `yaml:"Port"`
   AlsaDevice string `yaml:"AlsaDevice"`
   SpamInterval int `yaml:"SpamInterval"`
   //BottomSpeaks []string `yaml:"BottomSpeaks"`
   BottomSpeaks []BottomSpeak `yaml:"BottomSpeaks"`
   SignDevice string `yaml:"SignDevice"`
   NotifyChannel string `yaml:"NotifyChannel"`
   Mode string `yaml:"Mode"`
   SpeakHelp string `yaml:"SpeakHelp"`
}

func hello(w http.ResponseWriter, req *http.Request) {

    fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

    for name, headers := range req.Header {
        for _, h := range headers {
            fmt.Fprintf(w, "%v: %v\n", name, h)
        }
    }
}

var cfg SpeakConfig
var ch = make(chan SpeakRequest,1)
var persistText string = ""


func bottom(w http.ResponseWriter, req *http.Request) {
   if (req.Method == "POST") {
   		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(os.Stderr, "ParseForm() err: %v", err)
			return
		}

		text := req.PostFormValue("text")
		quickText := req.PostFormValue("quickText")
		audio := req.PostFormValue("audio")
    //log.Printf("Got text: %s\n",text)
    if ((cfg.SignDevice != "") && (quickText != "")) {
      alphasign(quickText,cfg.SignDevice)
    }
    if (cfg.SignDevice != "" && quickText == "" && text != "") {
      persistText = text
      alphasign(text,cfg.SignDevice)
    }

    if (len(audio) > 0) {
      ab, err := base64.URLEncoding.DecodeString(audio)
      if (err != nil) {
        //log.Printf("Base64 decode error %s",err)
      } else {
        //log.Printf("Got %d bytes PCM data",len(ab))
        if (cfg.AlsaDevice != "") {
          cmd:= exec.Command("aplay", "-D",cfg.AlsaDevice,"prompt.wav")
          cmd.Run()
          cmd= exec.Command("aplay", "-D",cfg.AlsaDevice,"-c","1","-f","S16_LE","-r","16000")
          cmd.Stdin = bytes.NewReader(ab)
          cmd.Run()
        }
      }
    }
    if (cfg.SignDevice != "" && quickText != "") {
      time.Sleep(10*time.Second)
      alphasign(persistText,cfg.SignDevice)
    }
    _ = audio
  }
}
var textout string
func slack(w http.ResponseWriter, req *http.Request) {

    for name, headers := range req.Header {
        for _, h := range headers {
            fmt.Fprintf(os.Stderr, "%v: %v\n", name, h)
        }
    }
   //buf := new(bytes.Buffer)
   //buf.ReadFrom(req.Body)
   //contents := buf.String()
   //fmt.Fprintf(os.Stderr,"BODY: %s\n",contents)
   //fmt.Fprintf(w,"Content-type: text/plain\n\n")

   /*
   var dat = make(map[string]string)
   fmt.Fprintf(w,contents)
   json.Unmarshal([]byte(contents),&dat)
   fmt.Fprintf(os.Stderr,"MARHSLLED %T %v+\n",dat,dat)
   for a,b := range dat {
	  fmt.Fprintf(os.Stderr,"GOT %v+ %v+\n",a,b)
   }
   */

   if (req.Method == "POST") {
   		if err := req.ParseForm(); err != nil {
			//log.Printf( "ParseForm() err: %v", err)
			return
		}

		token := req.PostFormValue("token")
		//log.Printf("TOKEN IS %s\n",token)
		if (token != cfg.Token) {
			fmt.Fprintf(os.Stderr,"Invalid token \"%s\"",token)
      fmt.Fprintf(os.Stderr,"Invalid Token")
      return
    }
		user_name := req.PostFormValue("user_name")
		//log.Printf("user_name IS %s\n",user_name)

		command := req.PostFormValue("command")
		//log.Printf("command IS %s\n",command)

		user_id := req.PostFormValue("user_id")
		//log.Printf("user_id IS %s\n",user_id)

		text := req.PostFormValue("text")
		//log.Printf("Text IS %s\n",text)



		//log.Printf("FORM IS %v\n",req.Form)
    log.Printf("%s:%s: %s",user_id,user_name,text)

    if (cfg.NotifyChannel != "") {
      post_slack(user_name,text)
    }


    if (cfg.Mode == "CGI") {
      // CGI is Depricated due to threading issues??
      /*
      id, _, _ := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
      if (id != 0) {
        fmt.Fprintf(w,"%s Anouncing: %s",user_name,text)
        fmt.Fprintf(os.Stderr,"Child exit\n")
        return
      }
      */
      fmt.Fprintf(os.Stderr,"Dispacthed speak\n")
      //speak(text) DONT SPEAK HERE! Will make CGI timeout. Thus function (HttpHandler) must exit first!
      textout = text
      fmt.Fprintf(os.Stderr,"Speak Done speak\n")
    } else if  ((command ==  "/speakhelp") || (command == "/helpspeak")) {
          fmt.Fprintf(w,"%s",cfg.SpeakHelp)
    } else {
      if (command[0] == '/') { command = command[1:]}
      speakreq := SpeakRequest{text,command}
      select {
        case ch <- speakreq:
          fmt.Fprintf(w,"%s Anouncing: %s",user_name,text)
        default:
          fmt.Fprintf(w,"Speakbot requires %d seconds between announcements. Please wait and retry",cfg.SpamInterval)
      }
    }

	}
}

// Add the waitgroup before calling!
func speak(text string,slashcmd string) {

    fmt.Fprintf(os.Stderr,"Speaking\n")
    awscfg, err := config.LoadDefaultConfig(context.TODO(),
      // Hard coded credentials.
      config.WithRegion("us-east-1"),
      config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
        Value: aws.Credentials{
          AccessKeyID: cfg.AccessKey, SecretAccessKey: cfg.SecretKey, SessionToken: "",
          Source: "example hard coded credentials",
        },
		}))
    if err != nil {
      fmt.Fprintf(os.Stderr,"AWS Polly Load error: %s\n",err)
      log.Fatal(err)
    }
    fmt.Fprintf(os.Stderr,"Poly1\n")
    p := polly.NewFromConfig(awscfg)

    fmt.Fprintf(os.Stderr,"Poly2\n")
    //log.Printf("CONFIG is %v+\n",cfg)
    input := &polly.DescribeVoicesInput{LanguageCode: "en-US"}
    _, err = p.DescribeVoices(context.TODO(),input)
    //log.Printf("Poly is %T %v+\n",p,p)
    //log.Printf("DescribeVoices %v is %T %v+\n",err,resp,resp)

    fmt.Fprintf(os.Stderr,"Poly3\n")
    t := fmt.Sprintf("<speak><amazon:domain name=\"news\"><prosody volume=\"x-loud\" rate=\"slow\">%s</prosody></amazon:domain></speak>",text)
    ssin := &polly.SynthesizeSpeechInput{
      OutputFormat: pollytype.OutputFormatPcm,
      LanguageCode: pollytype.LanguageCodeEnUs,
      TextType: pollytype.TextTypeSsml,
      VoiceId: pollytype.VoiceIdJoanna,
      Engine: pollytype.EngineNeural,
      Text: &t}
    //log.Printf(t)
    fmt.Fprintf(os.Stderr,"Synthesizing\n")
    spout, errout := p.SynthesizeSpeech(context.TODO(),ssin)
    if (errout != nil) {
      fmt.Fprintf(os.Stderr,"AWS Polly error: %s\n",errout)
    }
    pcmdata := new(bytes.Buffer)
    _,err = pcmdata.ReadFrom(spout.AudioStream)
    //log.Printf("Poly Read  %d Bytes %s\n",pcmdata.Len())


    //os.WriteFile("data.pcm", pcmdata.Bytes(), 0644)
    spout.AudioStream.Close()

    fmt.Fprintf(os.Stderr,"Aplaying\n")
    if (cfg.AlsaDevice != "") {
      cmd:= exec.Command("aplay", "-D",cfg.AlsaDevice,"prompt.wav")
      cmd.Run()
      cmd= exec.Command("aplay", "-D",cfg.AlsaDevice,"-c","1","-f","S16_LE","-r","16000")
      cmd.Stdin = bytes.NewReader(pcmdata.Bytes())
      cmd.Run()
    }

    fmt.Fprintf(os.Stderr,"Alphasigning\n")
    if (cfg.SignDevice != "") {
      alphasign(text,cfg.SignDevice)
    }
    fmt.Fprintf(os.Stderr,"Bottomspeaking\n")
    for _,bs := range cfg.BottomSpeaks {
        fmt.Fprintf(os.Stderr, "Match attempt: \"%s\"\n",bs)
        if (slices.Contains(bs.Commands,slashcmd)) {
          fmt.Fprintf(os.Stderr, "Dispatch to Bottom: \"%s\"\n",bs)
          response, err := http.PostForm(bs.URL, url.Values{
          "text": {text},
          "audio": {base64.URLEncoding.EncodeToString(pcmdata.Bytes())}})
          if (err != nil) {
            fmt.Fprintf(os.Stderr, "Bottom response is %v %s\n",response,err)
          }
        }
    }
     fmt.Fprintf(os.Stderr,"Speak is done\n")
}

func speaker() {
  for {
    speachreq := <- ch
    fmt.Fprintf(os.Stderr,"Got from Speaker: %s\n",speachreq.Text)
    //log.Println("Speaker got text",st)
    if (speachreq.Text != "") {
      ch <- SpeakRequest{"",""}
      speak(speachreq.Text,speachreq.SlashCommand)
      time.Sleep(time.Duration(cfg.SpamInterval) * time.Second)
    }
    //log.Println("Speaker loop done",st)
  }
}


func main() {


    f, err := os.Open("speak.cfg")
    decoder := yaml.NewDecoder(f)
    err = decoder.Decode(&cfg)
    if (err != nil) {
      log.Fatal("Config Decode error: ",err)
    }

    fmt.Fprintf(os.Stderr,"Speak log is %v\n",cfg)
    if (cfg.Mode == "CGI") {
      // CGI is Depricated due to threading issues??
      fmt.Fprintf(os.Stderr,"CGI Serving\n")
      err := cgi.Serve(http.HandlerFunc(slack))
      fmt.Fprintf(os.Stderr,"CGI Served\n")

      // NOW we can call "speak"!?
      //os.Stdout.Close()
      //os.Stdin.Close()
      /*
      id, _, _ := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
      if (id == 0) {
        fmt.Fprintf(os.Stderr,"Speak from Fork\n")
        speak(textout)
        fmt.Fprintf(os.Stderr,"Forked speak done\n")
      }
      */
        if err != nil {
          fmt.Printf("Status:%d %s\r\n", 500, "Cannot get request")
          fmt.Printf("Content-Type: text/plain\r\n")
          fmt.Printf("\r\n")
          fmt.Printf("%s\r\n", "Cannot get request")
            return
        }
        speak(textout,"")

        // Use req to handle request
      fmt.Fprintf(os.Stderr,"CGI Wait done\n")
      return
    }
    
    //log.Printf("Backends",cfg.BottomSpeaks)
    //awscfg, err := config.LoadDefaultConfig(context.TODO())
    http.HandleFunc("/hello", hello)
    http.HandleFunc("/headers", headers)
    http.HandleFunc("/slack", slack)
    http.HandleFunc("/bottom", bottom)

    log.Println("Listening Port",cfg.Port)
    speak("Speakbot active","/pa")
    go speaker()
    http.ListenAndServe(fmt.Sprintf(":%d",cfg.Port), nil)
}
