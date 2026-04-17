set dotenv-load := true
# set windows-shell := ["C:\\Program Files\\Git\\bin\\sh.exe","-c"]
set shell := ["powershell.exe", "-c"]

env name:
	ln -sf .env.{{name}} .env

start:
	go run cmd/service/main.go app start

test:
	go test ./...

lint:
	golangci-lint run ./... --fix
