clean:
	rm -rf ./bin

build: clean
	go mod download
	CGO_ENABLED=0 go build -o ./bin/certification-preparation-bot ./cmd/telegram/main.go

build-lambda-arm: clean
	go mod download
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/certification-preparation-bot-lambda-arm
	cp ./bin/certification-preparation-bot-lambda-arm ./bin/bootstrap
	zip -j ./bin/certification-preparation-bot-lambda-arm.zip ./bin/bootstrap