.PHONY: test run

test:
	godotenv -f ./.env ginkgo -v -r --cover

run:
	go run main.go --tlscert certs/server.crt --tlskey certs/server.key
