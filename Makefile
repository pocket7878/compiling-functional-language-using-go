fngo: main.go ast.go parser.go
	go build

parser.go: parser.y
	go generate

run: fngo
	cat ./sample.text | ./fngo