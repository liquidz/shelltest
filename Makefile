.PHONEY: clean test install dockertest

export PATH := .:$(PATH)
TARGET = shelltest

$(TARGET): *.go
	go build -o $(TARGET)

clean:
	\rm -f $(TARGET)

test: $(TARGET)
	go test && \
		./$(TARGET) example/hello_expected.txt && \
		./$(TARGET) example/regexp_expected.txt && \
		./$(TARGET) example/fail_expected.txt

install: $(TARGET)
	/bin/cp -pf $(TARGET) $(GOPATH)/bin

dockertest:
	\rm -f $(TARGET)
	docker run --rm -v `pwd`:/usr/src/myapp -w /usr/src/myapp golang make test
