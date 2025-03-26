# Makefile for festivals-gateway

VERSION=development
DATE=$(shell date +"%d-%m-%Y-%H-%M")
REF=refs/tags/development
DEV_PATH_MAC=$(shell echo ~/Library/Containers/org.festivalsapp.project)
export

build:
	go build -ldflags="-X 'github.com/Festivals-App/festivals-gateway/server/status.ServerVersion=$(VERSION)' -X 'github.com/Festivals-App/festivals-gateway/server/status.BuildTime=$(DATE)' -X 'github.com/Festivals-App/festivals-gateway/server/status.GitRef=$(REF)'" -o festivals-gateway main.go

install:
	mkdir -p $(DEV_PATH_MAC)/usr/local/bin
	mkdir -p $(DEV_PATH_MAC)/etc
	mkdir -p $(DEV_PATH_MAC)/var/log
	mkdir -p $(DEV_PATH_MAC)/usr/local/festivals-gateway

	cp operation/local/ca.crt  $(DEV_PATH_MAC)/usr/local/festivals-gateway/ca.crt
	cp operation/local/server.crt  $(DEV_PATH_MAC)/usr/local/festivals-gateway/server.crt
	cp operation/local/server.key  $(DEV_PATH_MAC)/usr/local/festivals-gateway/server.key
	cp festivals-gateway $(DEV_PATH_MAC)/usr/local/bin/festivals-gateway
	chmod +x $(DEV_PATH_MAC)/usr/local/bin/festivals-gateway
	cp operation/local/config_template_dev.toml $(DEV_PATH_MAC)/etc/festivals-gateway.conf

run:
	./festivals-gateway --container="$(DEV_PATH_MAC)"

test:
	go test ./server/loadbalancer

clean:
	rm -r festivals-gateway