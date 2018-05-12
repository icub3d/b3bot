GOBINARY=b3bot
GOSOURCES := $(shell find -name '*.go')

all: $(GOBINARY)

local: $(GOBINARY)
	./b3bot

$(GOBINARY): $(GOSOURCES)
	CGO_ENABLED=0 go build -a -installsuffix cgo -o $(GOBINARY) .

build: $(GOBINARY) Dockerfile
	docker build -t b3bot:latest .

run: build
	docker run -it b3bot
