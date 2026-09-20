FROM golang:1.25-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /service .
FROM alpine:3.22
COPY --from=build /service /service
EXPOSE 8080
CMD ["/service"]
