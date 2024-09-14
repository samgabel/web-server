.SILENT:

run:
	go run .

debug:
	go run . --debug

test:
	($(MAKE) wait-for-server && \
		curl -s -X POST localhost:8080/api/users --data '{"email":"test@testdomain.com", "password":"123456"}' > /dev/null 2>&1 && \
		curl -s -X POST localhost:8080/api/login --data '{"email":"test@testdomain.com", "password":"123456"}' | jq .refresh_token | \
		xargs -I _ curl -s -X POST localhost:8080/api/refresh --header 'Authorization: Bearer _' | jq .token | \
		xargs -I _ curl -s -X POST localhost:8080/api/chirps --header 'Authorization: Bearer _' --data '{"body":"hello, world!"}' > /dev/null 2>&1 && \
		curl -s -X POST localhost:8080/api/users --data '{"email":"another@testdomain.com", "password":"654321"}' > /dev/null 2>&1 && \
		curl -s -X POST localhost:8080/api/login --data '{"email":"another@testdomain.com", "password":"654321"}' | jq .refresh_token | \
		xargs -I _ curl -s -X POST localhost:8080/api/refresh --header 'Authorization: Bearer _' | jq .token | \
		xargs -I _ curl -s -X POST localhost:8080/api/chirps --header 'Authorization: Bearer _' --data '{"body":"the quick brown fox"}' > /dev/null 2>&1) &\
	go run . --debug


.PHONY: wait-for-server
wait-for-server:
	@echo "Waiting for server to start..."
	@until curl -s http://127.0.0.1:8080 >/dev/null 2>&1; do \
		sleep 1; \
	done
