.PHONY:build run dev stop redis mysql test
build:
	go build -o app main.go

run: build
	./app

dev: build
	nohup ./app > app.log 2>&1 & echo $$! > app.pid

stop:
	@if [ -f app.pid ] && [ -s app.pid ]; then \
		kill `cat app.pid` || true; \
		rm -f app.pid; \
	else \
		echo "no pid file"; \
	fi

restart: 
	@make stop 
	@makedev

redis:
	docker compose exec redis redis-cli

mysql:
	docker compose exec mysql  mysql -uroot -p

test:
	@bash scripts/api.sh