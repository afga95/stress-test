# Usar imagem oficial do Go como base para build
FROM golang:1.21-alpine AS builder

# Definir diretório de trabalho
WORKDIR /app

# Copiar o código fonte
COPY main.go .

# Inicializar módulo Go
RUN go mod init load-test-cli

# Compilar a aplicação
RUN go build -o load-test-cli main.go

# Usar imagem alpine minimalista para execução
FROM alpine:latest

# Instalar certificados SSL para requests HTTPS
RUN apk --no-cache add ca-certificates

# Criar diretório de trabalho
WORKDIR /root/

# Copiar o executável do estágio de build
COPY --from=builder /app/load-test-cli .

# Tornar o executável executável
RUN chmod +x load-test-cli

# Definir o executável como entrada padrão
ENTRYPOINT ["./load-test-cli"]