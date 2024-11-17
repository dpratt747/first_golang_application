FROM golang:latest AS build

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./main.go

# RUN ls -la .

# ----
FROM alpine:latest 

WORKDIR /root/

RUN apk add --no-cache libc6-compat  

# RUN ls -la .

COPY --from=build /app/main .
COPY --from=build /app/README.md .
COPY --from=build /app/.env .

RUN chmod +x /root/main

EXPOSE 8080 

CMD ["./main"]