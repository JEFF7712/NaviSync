.PHONY: all build clean install

TARGET := plugin.wasm
ARCHIVE := navisync.ndp

all: build

build:
	tinygo build -o $(TARGET) -target wasip1 .
	zip -j $(ARCHIVE) manifest.json $(TARGET)

clean:
	rm -f $(TARGET) $(ARCHIVE)

install: build
	@echo "To install, copy $(ARCHIVE) to your Navidrome plugins directory."
