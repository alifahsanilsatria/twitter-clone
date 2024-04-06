run-app:
	docker compose build && docker compose up -d
	docker logs -f backend >> /var/log/api-golang.log &

stop-app:
	docker compose down

run-debug-app:
	cd debugmode; docker compose build && docker compose up -d

stop-debug-app:
	cd debugmode; docker compose down