
all:
	go build -o ah-orchestrator

all-arm64:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o ah-orchestrator

run:
	./ah-orchestrator

clean:
	rm -rf ah-orchestrator
