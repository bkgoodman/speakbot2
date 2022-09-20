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
*/

import (
   "fmt"
   "net/http"
   "bytes"
   "context"
   "os"
   "os/exec"
   //"encoding/json"
   "gopkg.in/yaml.v2"

  "log"
  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/aws/aws-sdk-go-v2/config"
  "github.com/aws/aws-sdk-go-v2/service/polly"
  pollytype "github.com/aws/aws-sdk-go-v2/service/polly/types" 
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type SpeakConfig struct {
   SecretKey string `yaml:"SecretKey"`
   AccessKey string `yaml:"AccessKey"`
   BotToken string `yaml:"BotToken"`
   VerificationKey string `yaml:"VerificationKey"`
   Port int `yaml:"Port"`
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
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		user_name := req.PostFormValue("user_name")
		fmt.Fprintf(w,"user_name IS %s\n",user_name)

		user_id := req.PostFormValue("user_id")
		fmt.Fprintf(w,"user_id IS %s\n",user_id)

		text := req.PostFormValue("text")
		fmt.Fprintf(w,"test IS %s\n",text)

		token := req.PostFormValue("token")
		fmt.Fprintf(w,"TOKEN IS %s\n",token)


		fmt.Fprintf(os.Stderr,"FORM IS %v\n",req.Form)

	}
}

func speak(text string) {

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
      log.Fatal(err)
    }
    p := polly.NewFromConfig(awscfg)

    log.Printf("CONFIG is %v+\n",cfg)
    input := &polly.DescribeVoicesInput{LanguageCode: "en-US"}
    resp, err := p.DescribeVoices(context.TODO(),input)
    log.Printf("Poly is %T %v+\n",p,p)
    log.Printf("DescribeVoices %v is %T %v+\n",err,resp,resp)

    t := fmt.Sprintf("<speak><amazon:domain name=\"news\"><prosody volume=\"x-loud\" rate=\"slow\">%s</prosody></amazon:domain></speak>",text)
    ssin := &polly.SynthesizeSpeechInput{
      OutputFormat: pollytype.OutputFormatPcm,
      LanguageCode: pollytype.LanguageCodeEnUs,
      TextType: pollytype.TextTypeSsml,
      VoiceId: pollytype.VoiceIdJoanna,
      Engine: pollytype.EngineNeural,
      Text: &t}
    log.Printf(t)
    spout, errout := p.SynthesizeSpeech(context.TODO(),ssin)
    log.Printf("Poly Speak err %s\n",errout)
    log.Printf("Poly Speak  %s\n",*spout.ContentType)
    pcmdata := new(bytes.Buffer)
    _,err = pcmdata.ReadFrom(spout.AudioStream)
    log.Printf("Poly Read  %d Bytes %s\n",pcmdata.Len())


    os.WriteFile("data.pcm", pcmdata.Bytes(), 0644)
    spout.AudioStream.Close()
    cmd:= exec.Command("aplay", "-D","sysdefault:CARD=PCH","-c","1","-f","S16_LE","-r","16000")
    cmd.Stdin = bytes.NewReader(pcmdata.Bytes())
    cmd.Run()


}
func main() {

    f, err := os.Open("speak.cfg")
    decoder := yaml.NewDecoder(f)
    err = decoder.Decode(&cfg)
    if (err != nil) {
      log.Fatal("Config Decode error: ",err)
    }
    //awscfg, err := config.LoadDefaultConfig(context.TODO())
    http.HandleFunc("/hello", hello)
    http.HandleFunc("/headers", headers)
    http.HandleFunc("/slack", slack)

    log.Println("Listening Port",cfg.Port)
    speak("This is a test of the emergency broadcasting system. In the event of an actual emergency, you should put your head between your legs and kiss your ass goodbye.")
    http.ListenAndServe(fmt.Sprintf(":%d",cfg.Port), nil)
}
