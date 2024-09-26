BUILD_PATH=./bin/go_gof

compile:
	go build -o $(BUILD_PATH) main.go

run:
	$(BUILD_PATH)  -screenWidth=$(WIDTH) -screenHeight=$(HEIGHT) -filePath=$(FILE)

clean:
	rm -f $(BUILD_PATH)

all: clean compile run
