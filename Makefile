GH_USER ?= josegonzalez
NAME = sshd-config
HARDWARE = $(shell uname -m)
VERSION ?= 0.7.0

build: clean $(NAME)
	mkdir -p build/linux  && GOOS=linux  go build -ldflags "-X main.Version=$(VERSION)" -a -o build/linux/$(NAME)
	mkdir -p build/darwin && GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -a -o build/darwin/$(NAME)

clean:
	rm -rf build/* $(NAME)

run: $(NAME)
	./$(NAME)

$(NAME):
	go build -ldflags "-X main.Version=$(VERSION)"

release: build
	rm -rf release && mkdir release
	tar -zcf release/$(NAME)_$(VERSION)_linux_$(HARDWARE).tgz -C build/linux $(NAME)
	tar -zcf release/$(NAME)_$(VERSION)_darwin_$(HARDWARE).tgz -C build/darwin $(NAME)
	gh-release create $(GH_USER)/$(NAME) $(VERSION) $(shell git rev-parse --abbrev-ref HEAD)

.PHONY: build clean deps release run
