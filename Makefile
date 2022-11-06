.PHONY: all
all: speakbot speakbotrpi 

speakbot: speak.go
	go build -o speakbot
	
speakbotrpi: speak.go
	GOARCH=arm go build -o speakbotrpi


.PHONY: clean
clean:
	rm -rf speakbot speakbotrpi *~

.PHONY: updatelocal
updatelocal: speakbot
	systemctl stop bkgspeakbot.service
	cp speakbot /root/speakbot_live/speakbot
	systemctl start bkgspeakbot.service

.PHONY: remoteinstall
remoteinstall: speakbotrpi
	scp speakbotrpi bkg@speakbot:speakbot
	scp speakbotrpi bkg@garagebot:speakbot
