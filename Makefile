run:
	DOCKER_HOST_IP=$(shell ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+') docker-compose up --build