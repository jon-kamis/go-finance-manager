FROM golang:1.21.3 as builder

WORKDIR /app/finance-manager-backend

COPY go.mod go.sum ./
RUN go mod download
COPY cmd/** ./cmd/finance-mngr
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/finance-mngr

EXPOSE 8080

FROM golang:1.21.3
COPY --from=builder /app/finance-manager-backend/finance-mngr /finance-mngr
CMD ["/finance-mngr"]