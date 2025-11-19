FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN GOPROXY=https://goproxy.io,direct go mod download

COPY . .

# -ldflags "-s -w" reduces binary size by stripping debug information and symbol table
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main -ldflags "-s -w" .


FROM alpine:3.21

RUN apk update && apk upgrade && \
    apk add --no-cache bash make

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3002

CMD [ "./main", "runserver" ]