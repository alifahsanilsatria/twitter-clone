run-app:
	docker compose build && docker compose up -d
	docker logs -f backend >> /var/log/api-golang.log &

stop-app:
	docker compose down