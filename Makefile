.PHONEY: rebuild clean test install dockertest

export PATH := .:$(PATH)
TARGET = shelltest

$(TARGET): **/*.go
	go build -o $(TARGET)

rebuild: clean $(TARGET)

clean:
	\rm -f $(TARGET)

test: $(TARGET)
	go test ./... && \
		./$(TARGET) example/hello_expected.txt && \
		./$(TARGET) example/regexp_expected.txt && \
		./$(TARGET) example/fail_expected.txt && \
		./$(TARGET) example/fail_expected_tap.txt && \
		./$(TARGET) example/auto_assert_expected.txt && \
		./$(TARGET) example/auto_assert_no_expected.txt && \
		./$(TARGET) example/require_expected.txt

install: $(TARGET)
	/bin/cp -pf $(TARGET) $(GOPATH)/bin

dockertest:
	\rm -f $(TARGET)
	docker run --rm -v `pwd`:/usr/src/myapp -v ${GOPATH}:/go -w /usr/src/myapp golang make test
	\rm -f $(TARGET)
