.phony:
all: doorentry

.phony:
run: doorentry
	./doorentry 13

doorentry: doorentry.go
	go build -o doorentry

.PHONY: remoteinstall
remoteinstall: doorentry

	scp doorentry bkg@auth:doorentry
	ssh bkg@auth sudo /usr/bin/cp /home/bkg/doorentry /var/www/authbackend/doorentry
