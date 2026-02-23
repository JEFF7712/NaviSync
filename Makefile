.PHONY: all build clean install

TARGET := plugin.wasm
ARCHIVE := navisync.ndp

all: build

build:
	tinygo build -opt=2 -scheduler=none -no-debug -gc=leaking -o $(TARGET) -target wasip1 -buildmode=c-shared .
	zip -j $(ARCHIVE) manifest.json $(TARGET)

clean:
	rm -f $(TARGET) $(ARCHIVE)

install: build
	@echo "To install, copy $(ARCHIVE) to your Navidrome plugins directory."
