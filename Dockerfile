FROM golang:1.21-alpine
RUN apk add --no-cache --update gcc musl-dev
WORKDIR /app
COPY . .

RUN CGO_ENABLED=1 go build -o /backend

CMD [ "/backend" ]