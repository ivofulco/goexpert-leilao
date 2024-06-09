infra: docker-build
	docker compose up -d;

docker-build:
	docker compose build;

docker-down:
	docker compose down;	

#create-auction:
#   curl -X POST http://localhost:8080/auction \
#   -H "Content-Type: application/json" \
#   -d '{
#        "product_name": "Piano",
#        "category": "Music",
#        "description": "Classical Musical Instrument",
#        "condition": 1
#       }'   

get:
	@sleep 3s ;
	curl -X GET http://localhost:8000/order

desafio: infra



docker-clean:
	sudo docker rm -f $(docker ps -a -q)
	sudo docker rmi -f $(docker images -q)
	sudo docker system prune -af
	sudo docker volume prune -f