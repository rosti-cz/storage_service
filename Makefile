HOST=rosti-jupiter

build:
	CGO_ENABLED=0 go build -o storage_service *.go

deploy: build
	scp storage_service ${HOST}:/usr/local/bin/storage_service.tmp
	ssh ${HOST} mv /usr/local/bin/storage_service.tmp /usr/local/bin/storage_service
	ssh ${HOST} systemctl restart rosti_storageservice
