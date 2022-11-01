default:
.SILENT:

IMAGE_NAME := "g5-backend-service"

build:
	@docker build -t $(IMAGE_NAME) .

run:
	@docker run --rm -p 8000:8000 $(IMAGE_NAME)
