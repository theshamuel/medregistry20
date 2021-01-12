OS=linux
ARCH=amd64

image-prod:
	docker build -t theshamuel/medregestry20 .

image-dev:
	docker build -t theshamuel/medregestry20 --build-arg SKIP_TESTS=true .

deploy:
	docker-compose up -d

.PHONY: image
