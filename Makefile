
# build-db-image:
# 	docker build -t kids-bank-sqlite:latest -f ./Dockerfile .

# run: build-db-image
# 	docker-compose up --build -d
# 	@echo "Waiting for sqlite container to be healthy..."
# 	@while [ "$$(docker inspect --format='{{json .State.Health.Status}}' kids-bank-sqlite | jq -r .)" != "healthy" ]; do \
# 		echo "Waiting for kids-bank-sqlite to be healthy..."; \
# 		sleep 1; \
# 	done
# 	@echo "sqlite container is healthy"
# 	@sleep 0.5

run:
	@echo "Starting server..."
	KB_SERVER="yas" go tool air