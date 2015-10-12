.PHONEY: test install

TARGET = shelltest

$(TARGET): *.go
	go build

test: $(TARGET)
	go test && ./$(TARGET) example/shelltest.txt

install: $(TARGET)
	/bin/cp -pf $(TARGET) $(GOPATH)/bin
