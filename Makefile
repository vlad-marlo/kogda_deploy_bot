.PHONY: build
build:
	go build -o telebot ./cmd/telebot
.PHONY: dock
dock:
	docker buildx build --push --platform linux/amd64,linux/arm64 --file=infra/bot.dockerfile --tag="vladmarlo/kogda_deploy:latest" .

.PHONY: dock/push
dock/push: dock
	docker push vladmarlo/kogda_deploy:latest