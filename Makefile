include .env
export

.SILENT:
run: build-updater
	copy PoeBuyUpdater.exe .\cmd\app\PoeBuyUpdater.exe && \
	go run ./cmd/app/main.go

build: build-updater
	copy PoeBuyUpdater.exe .\cmd\app\PoeBuyUpdater.exe && \
	fyne package --src ./cmd/app --os windows --release && \
	move /Y .\cmd\app\PoeBuy.exe .

build-updater:
	go build -ldflags -H=windowsgui -o PoeBuyUpdater.exe ./cmd/updater/main.go