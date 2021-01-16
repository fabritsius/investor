include .env
export $(shell sed 's/=.*//' .env)

build:
	go build -o investor

dev:
	go run main.go