all: build run

build:
	go build -o ./build/go_gof .

run:
	./build/go_gof $(FILE)
