# Makefile to run the nationpulse project with just one command locally

run: start
	
# frontend:
# 	cd ./Frontendpnpm run dev

start: # run the instance of redis, postgres, kafka, prmetheus 
	docker compose up -d --build bff ingestion cronjob prometheus reporting kafka-1 kafka-2 kafka-3 kafka-init
	docker compose logs -f bff reporting ingestion cronjob prometheus

stop:
	docker compose down
