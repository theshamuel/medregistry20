OS=linux
ARCH=amd64
BUILDKIT_PROGRESS=plain
IMG_VER?=latest

image-prod:
	docker build -t theshamuel/medregapi-v2:${IMG_VER} .

image-dev:
	docker build -t theshamuel/medregapi-v2:${IMG_VER} --build-arg SKIP_TESTS=true .

deploy:
	docker-compose up -d

.PHONY: image-dev deploy
