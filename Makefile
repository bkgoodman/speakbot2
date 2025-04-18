.PHONY: all
all: speakbot speakbotrpi 

speakbot: speak.go
	go build -o speakbot
	
speakbotrpi: speak.go
	GOARCH=arm go build -o speakbotrpi


.PHONY: clean
clean:
	rm -rf speakbot speakbotrpi *~

.PHONY: remoteinstall
remoteinstall: speakbotrpi speakbot
	echo TODO - need to shutdown and restart for this to work
	ssh bkg@speakbot sudo systemctl stop bkgspeakbot.service
	scp speakbotrpi bkg@speakbot:speakbot
	ssh bkg@speakbot sudo systemctl start bkgspeakbot.service

	ssh bkg@garagebot sudo systemctl stop bkgspeakbot.service
	scp speakbotrpi bkg@garagebot:speakbot
	ssh bkg@garagebot sudo systemctl start bkgspeakbot.service

	ssh bkg@cgimisc sudo systemctl stop bkgspeakbot.service
	scp speakbot bkg@cgimisc:
	ssh bkg@cgimisc sudo cp speakbot /home/speakbot/speakbot_live/speakbot
	ssh bkg@cgimisc sudo systemctl start bkgspeakbot.service

	ssh bkg@speakbot2 sudo systemctl stop bkgspeakbot.service
	scp speakbotrpi bkg@speakbot2:speakbot
	ssh bkg@speakbot2 sudo systemctl start bkgspeakbot.service
