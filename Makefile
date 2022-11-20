bin/mqtt-to-google-form: cmd/mqtt-to-google-form/main.go cmd/mqtt-to-google-form/config.go
	  GOOS=linux GOARCH=arm go build -o $@ ./cmd/mqtt-to-google-form

install:
		cp ./bin/mqtt-to-google-form /usr/local/bin
