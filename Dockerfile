FROM golang:1.18-alpine

WORKDIR /app

# BUILD_TARGET: client or server
ARG BUILD_TARGET

COPY go.mod ./
RUN go mod download

COPY . ./

RUN go build -o ./bin/$BUILD_TARGET ./cmd/$BUILD_TARGET/main.go

RUN chmod +x ./bin/$BUILD_TARGET

EXPOSE 5555

ENV MAIN_BINARY=/app/bin/${BUILD_TARGET}

CMD ${MAIN_BINARY}