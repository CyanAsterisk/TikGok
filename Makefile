.PHONY: start
start:
	docker-compose up -d

.PHONY: stop
stop:
	docker-compose down

# run the interaction
.PHONY: interaction
interaction:
	go run ./server/cmd/interaction

# run the sociality
.PHONY: socaility
sociality:
	go run ./server/cmd/sociality

# run the user
.PHONY: user
user:
	go run ./server/cmd/user

# run the video
.PHONY: video
video:
	go run ./server/cmd/video

# run the api
.PHONY: api
api:
	go run ./server/cmd/api