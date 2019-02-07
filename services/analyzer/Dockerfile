FROM golang:1.11-alpine as build
RUN apk add --no-cache --update alpine-sdk git gcc
COPY . /build
WORKDIR /build
#RUN go get github.com/dghubble/sling github.com/gin-gonic/gin github.com/kelseyhightower/envconfig
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"'

FROM scratch
COPY --from=build /build/analyzer /
CMD ["/analyzer"]