.PHONEY: clean test install dockertest

TARGET = shelltest

$(TARGET): *.go
	go build -o $(TARGET)

clean:
	\rm -f $(TARGET)

test: $(TARGET)
	go test && env PATH=`pwd`:$(PATH) ./$(TARGET) example/shelltest.txt

install: $(TARGET)
	/bin/cp -pf $(TARGET) $(GOPATH)/bin

dockertest:
	\rm -f $(TARGET)
	docker run --rm -v `pwd`:/usr/src/myapp -w /usr/src/myapp golang make test
