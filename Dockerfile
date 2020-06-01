# build stage
FROM golang:1.13
ENV GO111MODULE=on
WORKDIR /app/go-mutant-detector
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o ./out/mutants-detector .


EXPOSE 8080

CMD ["./out/mutants-detector"]