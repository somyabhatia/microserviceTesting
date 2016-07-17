.PHONY: all image clean publish

IMAGE=rade/dummi

all: image

dummi: main.go
	go build -tags netgo

image: Dockerfile dummi
	docker build -t $(IMAGE) .

clean:
	rm -f dummi
	docker rmi -f $(IMAGE) 2>/dev/null || true

publish:
	docker push $(IMAGE)
