all:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o cloudcms .

docker:
	- docker image rm reg.urantiatech.com/cloudcms/cloudcms
	docker build -t reg.urantiatech.com/cloudcms/cloudcms .
