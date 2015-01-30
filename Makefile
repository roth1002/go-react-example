BIN      = $(GOPATH)/bin
NODE_BIN = ./node_modules/.bin
PID      = .pid
GO_FILES = $(filter-out bindata.go, $(wildcard *.go))
STATIC   = $(filter-out static/bundle.js, $(shell find static -type f))
APP      = $(shell find app -type f)
BUNDLE   = static/bundle.js

build: bindata.go

serve:
	@make restart
	@fswatch $(GO_FILES) $(STATIC) $(APP) | xargs -n1 -I{}  make restart || make kill

kill:
	@kill `cat $(PID)` || true

bindata.go: $(STATIC) $(BUNDLE)
	$(BIN)/go-bindata -pkg=main -prefix=static -o=$@ static/...

$(BUNDLE): $(BUNDLE:static%=app%)
	@mkdir -p $(@D)
	$(NODE_BIN)/browserify -e $< > $@

restart: build
	@make kill
	go build -o server
	./server & echo $$! > $(PID)

.PHONY: build serve restart kill stuff
