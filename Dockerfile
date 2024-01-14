FROM golang:1.21-alpine3.17 AS builder
WORKDIR /go/app
COPY . ./
RUN go mod download
RUN mkdir bin
RUN cd cmd/app/ && go build -o ../../bin/go-portal

FROM node:current-alpine3.18 AS asset-builder
WORKDIR /frontend
COPY --from=builder /go/app/frontend .
RUN npm i
RUN npm run build

FROM golang:1.21-alpine3.17
WORKDIR /go
EXPOSE 8888
COPY --from=builder /go/app/bin/go-portal ./
COPY --from=asset-builder /frontend/dist ./frontend/dist
ENTRYPOINT ["./go-portal"]