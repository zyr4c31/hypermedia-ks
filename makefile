dev:
	templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."

run:
	templ generate
	go run .
