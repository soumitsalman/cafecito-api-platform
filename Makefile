.PHONY: gateway-dev beans-dev espresso-dev docker-up docker-beans docker-espresso

gateway-dev:
	npm run dev

beans-dev:
	cd services/beans && make run

espresso-dev:
	cd services/espresso && make run

docker-up:
	docker compose up --build

docker-beans:
	docker compose up --build tei beansapi

docker-espresso:
	docker compose up --build tei espressoapi
