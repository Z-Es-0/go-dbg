build:
	go build -o bin\debug.exe main\main.go


run:
	go run main\main.go

test:
	./bin\debug.exe sever\test.exe
	
	
