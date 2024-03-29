default:
.SILENT:

IMAGE_NAME := "ghcr.io/game-ai-project/g5-backend-service:latest"

build:
	@docker rmi -f $(IMAGE_NAME) || true
	@docker build -t $(IMAGE_NAME) .

run:
	@docker run --rm -p 8000:8000 --env-file=.env $(IMAGE_NAME)

publish:
	@docker push $(IMAGE_NAME)
