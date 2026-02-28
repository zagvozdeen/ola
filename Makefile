.PHONY: dev test build deploy

dev:
	GOEXPERIMENT=jsonv2 go run cmd/main.go -config=config.local.yaml

test:
	GOEXPERIMENT=jsonv2 go test ./...

build:
	GOEXPERIMENT=jsonv2 GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ola cmd/main.go
	npm run build

deploy: build
	ssh deploy@178.212.13.87 "sudo systemctl stop ola.service || true && cd /var/www/olastudio-ekb.ru && rm -rf public templates ola && mkdir -p public templates && chgrp ola public templates && chmod 2755 public && chmod 2750 templates"
	scp ola deploy@178.212.13.87:/var/www/olastudio-ekb.ru
	scp config.production.yaml deploy@178.212.13.87:/var/www/olastudio-ekb.ru/config.yaml
	tar -czf /tmp/ola-public.tar.gz -C dist .
	scp /tmp/ola-public.tar.gz deploy@178.212.13.87:/tmp/ola-public.tar.gz
	ssh deploy@178.212.13.87 "mkdir -p /var/www/olastudio-ekb.ru/public && tar -xzf /tmp/ola-public.tar.gz -C /var/www/olastudio-ekb.ru/public && rm -f /tmp/ola-public.tar.gz"
	rm -f /tmp/ola-public.tar.gz
	scp -r templates/* deploy@178.212.13.87:/var/www/olastudio-ekb.ru/templates/
	ssh deploy@178.212.13.87 "chmod -R g+rX,o-rwx /var/www/olastudio-ekb.ru/templates"
#	scp deploy/nginx.conf deploy@178.212.13.87:/etc/nginx/sites-available/olastudio-ekb.ru
	ssh deploy@178.212.13.87 "sudo systemctl daemon-reload && sudo systemctl restart ola.service && sudo systemctl reload nginx"
	rm ola && rm -rf dist
