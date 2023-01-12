.PHONY: start
start:
	docker-compose up -d

.PHONY: stop
stop:
	docker-compose down