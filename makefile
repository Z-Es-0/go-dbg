build:
	go build -o bin\debug.exe sever\main.go

test:
	go build -o bin\debug.exe sever\main.go
	./bin/debug.exe sever/test.exe
	
	