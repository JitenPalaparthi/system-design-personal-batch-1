FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -o /app main.go
FROM alpine:3.22
COPY --from=build /app /app
EXPOSE 8080
CMD ["/app"]
