version="0.1.0"
version_file=VERSION
working_dir=$(shell pwd)
arch="armhf"
remote_host = "fh@cube.local"

clean:
	-rm  -f ./src/mill
	find . -name '.DS_Store' -type f -delete

build-go:
	cd ./src;go build -o mill service.go;cd ../

build-go-arm:
	cd ./src;GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o mill service.go;cd ../

build-go-amd:
	cd ./src;GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o mill service.go;cd ../


configure-arm:
	python ./scripts/config_env.py prod $(version) armhf

configure-amd64:
	python ./scripts/config_env.py prod $(version) amd64

package-tar:
	tar cvzf mill_$(version).tar.gz mill VERSION

clean-deb:
	find package/debian -name ".DS_Store" -delete
	find package/debian -name "delete_me" -delete
	find package/debian -name ".DS_Store" -delete
	find package/debian -name "delete_me" -delete

package-deb-doc:clean-deb
	@echo "Packaging application using Thingsplex debian package layout"
	chmod a+x package/debian/DEBIAN/*
	mkdir -p package/debian/var/log/thingsplex/mill package/debian/var/lib/futurehome/mill/data package debian/usr/bin
	mkdir -p package/build
	cp ./src/mill package/debian/opt/thingsplex/mill
	cp $(version_file) package/debian/var/lib/futurehome/mill
	docker run --rm -v ${working_dir}:/build -w /build --name debuild debian dpkg-deb --build package/debian
	@echo "Done"

package-docker-amd:build-go-amd
	cp ./src/mill package/docker/service
	cd ./package/docker;docker build -t mill .

deb-arm : clean configure-arm build-go-arm package-deb-doc
	@echo "Building Thingsplex ARM package"
	@mv package/debian.deb package/build/mill_$(version)_armhf.deb

deb-amd : configure-amd64 build-go-amd package-deb-doc
	@echo "Building Thingsplex AMD package"
	mv package/debian.deb mill_$(version)_amd64.deb

upload :
	@echo "Uploading the package to remote host"
	scp package/build/mill_$(version)_armhf.deb $(remote_host):~/

remote-install : upload
	@echo "Uploading and installing the package on remote host"
	ssh -t $(remote_host) "sudo dpkg -i mill_$(version)_armhf.deb"

deb-remote-install : deb-arm remote-install
	@echo "Package was built and installed on remote host"


run :
	cd ./src; go run service.go -c ../testdata;cd ../


.phony : clean
