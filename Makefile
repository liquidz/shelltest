.PHONEY: rebuild clean test install dockertest

export PATH := .:$(PATH)
TARGET = shelltest

$(TARGET): *.go **/*.go
	go build -o $(TARGET)

rebuild: clean $(TARGET)

clean:
	\rm -f $(TARGET)

test: $(TARGET)
	go test ./... && \
		./$(TARGET) check_ver.shelltest \
					example/hello_expected.txt  \
					example/regexp_expected.txt  \
					example/fail_expected.txt  \
					example/fail_expected_tap.txt  \
					example/auto_assert_expected.txt  \
					example/auto_assert_no_expected.txt  \
					example/require_expected.txt

install: $(TARGET)
	/bin/cp -pf $(TARGET) $(GOPATH)/bin

dockertest:
	\rm -f $(TARGET)
	docker run --rm -v `pwd`:/usr/src/myapp -v ${GOPATH}:/go -w /usr/src/myapp golang make test
	\rm -f $(TARGET)
