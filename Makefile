spell:
	npx --yes cspell "**/*" --exclude="**/target/**" --exclude="**/out/**"

run:
	go run main.go

build:
	go build -o main .
