.PHONY: all clean build package build-linux build-windows prepare-package package-linux package-windows

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

all: build package

build: build-linux build-windows

package: prepare-package package-linux package-windows

build-linux:
	@echo "Building ttm (linux)..."
	mkdir -p bin
	go build -o bin/ttm .

build-windows:
	@echo "Building ttm (windows)..."
	mkdir -p bin
	GOOS=windows GOARCH=amd64 go build -o bin/ttm.exe .

prepare-package:
	rm -rf dist
	mkdir -p dist

package-linux: build-linux
	@echo "Copying package structure..."
	mkdir -p dist
	cp -r build/linux/ttm dist/ttm
	mkdir -p dist/ttm/usr/local/bin
	cp bin/ttm dist/ttm/usr/local/bin/ttm
	@echo "Creating DEB package..."
	dpkg-deb --build dist/ttm dist/ttm.deb
	rm -rf dist/ttm
	@echo "DEB package created at dist/ttm.deb"

package-windows: build-windows
	@command -v makensis >/dev/null 2>&1 || { echo "makensis (NSIS) not found. Install 'nsis' (e.g. sudo apt install nsis)"; exit 1; }
	@echo "Creating Windows installer..."
	mkdir -p dist
	cp -r build/windows dist/
	cp bin/ttm.exe dist/windows/
	cd dist/windows && makensis -DVERSION="$(VERSION)" ttm_installer.nsi
	@if [ -f dist/windows/ttm_setup.exe ]; then mv dist/windows/ttm_setup.exe dist/ttm_setup.exe; fi
	rm -rf dist/windows
	@echo "Windows installer created at dist/ttm_setup.exe"

clean:
	rm -rf dist
	rm -rf bin