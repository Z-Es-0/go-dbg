build:
	go build -o bin\debug.exe main\main.go


run:
	go build -o bin\debug.exe main\main.go
	.\bin\debug.exe "E:\ZesOJ\sever\test.exe"

test:
	./bin\debug.exe sever\test.exe
	
	
