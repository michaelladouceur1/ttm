.PHONY: all clean build package build-linux build-windows prepare-package package-linux package-windows 

all: build package

build: build-linux build-windows

package: prepare-package package-linux package-windows

build-linux:
	@echo "Building ttm (linux)..."
	go build -o bin/ttm main.go

build-windows:
	@echo "Building ttm (windows)..."
	GOOS=windows GOARCH=amd64 go build -o bin/ttm.exe main.go

prepare-package:
	rm -rf dist
	mkdir -p dist

package-linux: build
	@echo "Copying package structure..."
	rm -rf dist
	mkdir -p dist
	cp -r build/linux/ttm dist/ttm
	cp bin/ttm dist/ttm/usr/local/bin/ttm
	@echo "Creating DEB package..."
	dpkg-deb --build dist/ttm dist/ttm.deb
	rm -rf dist/ttm
	@echo "DEB package created at dist/ttm.deb"

package-windows: build-windows
	@echo "Creating Windows installer..."
	cp -r build/windows dist/
	cp bin/ttm.exe dist/windows/
	makensis dist/windows/ttm_installer.nsi
	mv dist/windows/ttm_setup.exe dist/ttm_setup.exe
	rm -rf dist/windows
	@echo "Windows installer created at dist/ttm_setup.exe"

clean:
	rm -rf dist
	rm -rf bin