#  Speakbot

## Modes
Speakbot runs in two different *modes*

- `CGI` App is called via a web server as a CGI when a request comes in. (*DEPRICATED*? Threading issues???)
- `Daemon` run the app and it will start it's own daemon, awaiting requests on the speified `Port`

CGI mode does NOT WORK properly and will probably never. This is because the upstream CGI provider (Apache) waits for the process to _terminate_ before closing the connection. This inherently means the entire process here has to have been completed. Slack on the other hand wants a response within 3 seconds, and if it doesn't get one returns a timeout. Thus, this needs to run as a standalone daemon which can respond to the incomming Slack request, _close_ the connection, then continue to handle the speech in another thread/goroutine. (I tried forks and subprocesses, this inherently couldn't fix this problem, or created others)


## Front and Back Ends

Speakbot can run as a "front end" or a "back end" (or both!)

### Front End

"Front end" means it's talking directly to slack and taking incoming requests. This means it Slack command will reach out to it's URL (weither running as a Daemon or CGI, per above) - and this code will reach out to Amazon to synthesize the audio.

From here: 

- If you have configured one or more `BottomSpeak` URLs - it will send the text and the already syntehsized audio to a "Back End" at the specified URL

### Back End

"Back End" means it will be accepting audio which has already been synthesized from a "front end". Thus:

- If you have configured an `AlsaDevice` - it will use that ALSA sound device to play the audio. Use the `aplay -L` command for a list of device. Not all will work you probably want one of the generic system default devices. (i.e. some Low-level devices will ONLY work with certian audio formats, whereas the "higher" level ones will do whatver conversion is necessary. If something doesn't work - try a different one)
- If you have configured `SignDevice` (the name of the serial device for a AlphaSign display like `/dev/ttyUSB0`) - it will send the text to that device


### Both

In theory - you could have a single node which acts as both a "Front End" and a "Back End", by recieving stuff directly from slack, synthesizing, playing and sending to sign.

Even more theoretical - you could do all this _and_ continue to send it off to yet another backend.

# Configuration
The `SecretKey` and `AccessKey` come from Amazon - and must contain AIM credentials which permit you to use Polly audio synthesis.

BotToken isn't used right now - but was intented to hold Slack "bot" credentials if you needed to make a Slack API call.

`Token` is your Slack Applications "Verification Token" (listed under your apps "Basic Information". If you don't put it here - speakbot will reject messages from Slack

## Slack Setup

Create your bot in app, and configure as a "Slash Command". Just give it the URL where Slack can reach a running Speechbot

`NotifyChannel` is the chanel identifier if you want speak commands logged to a Slack channel


# Doorentry

Note the `doorentry` directory. This is a separate script run on `auth` server via `mqtt_daemon` when a new person enters building. `Doorentry` needs to be installed alongside `mqtt_daemon` and has it's own config file. It does n't go through the main CGI hander above, but rather speaks out directly to the individial backends, delivering welcome message and member audio.
