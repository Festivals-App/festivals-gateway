# Makefile for festivals-gateway

VERSION=development
DATE=$(shell date +"%d-%m-%Y-%H-%M")
REF=refs/tags/development
export

build:
	go build -v -ldflags="-X 'github.com/Festivals-App/festivals-gateway/server/status.ServerVersion=$(VERSION)' -X 'github.com/Festivals-App/festivals-gateway/server/status.BuildTime=$(DATE)' -X 'github.com/Festivals-App/festivals-gateway/server/status.GitRef=$(REF)'" -o festivals-gateway main.go

install:
	cp festivals-gateway /usr/local/bin/festivals-gateway
	cp config_template.toml /etc/festivals-gateway.conf
	cp operation/service_template.service /etc/systemd/system/festivals-gateway.service

update:
	systemctl stop festivals-gateway
	cp festivals-gateway /usr/local/bin/festivals-gateway
	systemctl start festivals-gateway

uninstall:
	systemctl stop festivals-gateway
	rm /usr/local/bin/festivals-gateway
	rm /etc/festivals-gateway.conf
	rm /etc/systemd/system/festivals-gateway.service

run:
	./festivals-gateway

test:
	go test ./server/loadbalancer

stop:
	killall festivals-gateway

clean:
	rm -r festivals-gateway