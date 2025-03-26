FROM golang:1.23-alpine as build

WORKDIR /src/allocator
ADD . /src/allocator
RUN CGO_ENABLED=0 go build -o /bin/allocator .

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build /bin/allocator /app/allocator
EXPOSE 8000
CMD ["/app/allocator"]
