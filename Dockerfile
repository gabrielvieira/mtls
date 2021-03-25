FROM golang:1.15.5

WORKDIR /build

COPY . .

EXPOSE 8080

RUN  go build -o app server/server.go

CMD [ "/build/app" ]