
all: generate build deploy

run:
	go run -tags=develop ./server

generate:
	go generate ./server

build:
	GOOS=linux GOARCH=arm GOARM=6 go build -o rpicon ./server

deploy:
	scp rpicon procon.py pi@procon.local:/home/pi/
