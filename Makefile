TAG=$(shell git rev-parse --short ${GITHUB_SHA})$(and $(shell git status -s),-dirty)

build-allocator:
	docker buildx build --platform linux/amd64 -f docker/Dockerfile -t murmurations/allocator:latest .

tag-allocator: build-allocator
	docker tag murmurations/allocator murmurations/allocator:$(TAG)

push-allocator: tag-allocator
	docker push murmurations/allocator:latest
	docker push murmurations/allocator:$(TAG)

deploy-allocator:
	helm upgrade murmurations-allocator ./murmurationsAllocator --set image=murmurations/allocator:$(TAG) --install --atomic