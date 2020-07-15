fngo: *.go
	go build

parser.go: parser.y
	go generate

run: fngo
	cat ./sample.text | ./fngo