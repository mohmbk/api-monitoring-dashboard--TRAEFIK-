FROM golang:alpine as builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build  -ldflags="-s -w" -o /src/backend .


FROM scratch 
COPY --from=builder /src/backend .
EXPOSE 8080 27017
ENTRYPOINT [ "/backend" ]
