.PHONY: dev build deploy

dev:
	GOEXPERIMENT=jsonv2 go run cmd/main.go -config=config.local.yaml

build:
	GOEXPERIMENT=jsonv2 GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ola cmd/main.go
	npm run build

deploy: build
	ssh root@185.221.214.4 "service ola stop && cd /var/www/ola.creavo.ru && rm -rf public && rm -rf templates && rm ola || true"
	scp ola root@185.221.214.4:/var/www/ola.creavo.ru
	scp config.production.yaml root@185.221.214.4:/var/www/ola.creavo.ru/config.yaml
	tar -czf /tmp/ola-public.tar.gz -C dist .
	scp /tmp/ola-public.tar.gz root@185.221.214.4:/tmp/ola-public.tar.gz
	ssh root@185.221.214.4 "mkdir -p /var/www/ola.creavo.ru/public && tar -xzf /tmp/ola-public.tar.gz -C /var/www/ola.creavo.ru/public && rm -f /tmp/ola-public.tar.gz"
	rm -f /tmp/ola-public.tar.gz
#	tar -czf /tmp/ola-templates.tar.gz -C templates .
#	scp /tmp/ola-templates.tar.gz root@185.221.214.4:/tmp/ola-templates.tar.gz
#	ssh root@185.221.214.4 "rm -rf /var/www/ola.creavo.ru/templates && mkdir -p /var/www/ola.creavo.ru/templates && tar -xzf /tmp/ola-templates.tar.gz -C /var/www/ola.creavo.ru/templates && rm -f /tmp/ola-templates.tar.gz"
#	rm -f /tmp/ola-templates.tar.gz
	scp -r templates root@185.221.214.4:/var/www/ola.creavo.ru/templates
	scp deploy/nginx.conf root@185.221.214.4:/etc/nginx/sites-available/ola.creavo.ru
#	ssh root@185.221.214.4 "ln -s /etc/nginx/sites-available/ola.creavo.ru /etc/nginx/sites-enabled/ola.creavo.ru"
#	scp deploy/ola.service root@185.221.214.4:/etc/systemd/system
	ssh root@185.221.214.4 "systemctl daemon-reload && service ola restart && nginx -s reload"
#	rm ola && rm -rf dist
	rm ola
