FROM golang:1.21.3 as builder

WORKDIR /app/finance-manager-backend

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/go-finance-manager

EXPOSE 8080

FROM golang:1.21.3
COPY --from=builder /app/finance-manager-backend/cmd /finance-manager
COPY --from=builder /app/finance-manager-backend/docs /finance-manager
COPY --from=builder /app/finance-manager-backend/API.key /finance-manager
#CMD ["/finance-manager/cmd/go-finance-manager"]
ENTRYPOINT ["#!/bin/bash"]