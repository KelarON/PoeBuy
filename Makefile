# include .env
# export

.SILENT:
run:
	go run main.go

build:
	fyne package -os windows --release