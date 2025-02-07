dev:
	templ generate --watch --proxy="http:192.168.3.112:8080" --cmd="go run ."

run:
	templ generate
	go run .
