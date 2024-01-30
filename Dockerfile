FROM golang:alpine AS build

WORKDIR /build
ADD . /build

RUN go build .

FROM alpine:latest

# Export necessary port
EXPOSE 80
# Add  application
WORKDIR /dist
COPY --from=build /build/pagespeed-performance-svg /dist/main

ENTRYPOINT ["/dist/main"]