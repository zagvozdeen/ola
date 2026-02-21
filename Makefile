.PHONY: dev build deploy

dev:
	GOEXPERIMENT=jsonv2 go run cmd/main.go

build:
	GOEXPERIMENT=jsonv2 GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ola cmd/main.go
	npm run build

deploy: build
	ssh root@185.221.214.4 "service ola stop && cd /var/www/ola.creavo.ru && rm -rf public && rm ola || true"
	scp ola root@185.221.214.4:/var/www/ola.creavo.ru
	scp .env.production root@185.221.214.4:/var/www/ola.creavo.ru/.env
	scp -r dist root@185.221.214.4:/var/www/ola.creavo.ru/public
	scp deploy/nginx.conf root@185.221.214.4:/etc/nginx/sites-available/ola.creavo.ru
	ssh root@185.221.214.4 "ln -s /etc/nginx/sites-available/ola.creavo.ru /etc/nginx/sites-enabled/ola.creavo.ru"
	scp deploy/ola.service root@185.221.214.4:/etc/systemd/system
	ssh root@185.221.214.4 "systemctl daemon-reload && service ola restart && nginx -s reload"
#	rm ola && rm -rf dist