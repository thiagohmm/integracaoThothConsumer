# Estágio de construção (builder)
FROM golang:1.24.1-alpine AS builder

# Define o diretório de trabalho
WORKDIR /app

# Copia os arquivos de dependências (go.mod e go.sum)
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia todo o código-fonte para o diretório de trabalho
COPY . .

# Define o diretório de trabalho para o diretório do aplicativo
WORKDIR /app/cmd/app

# Configura variáveis de ambiente para a compilação
ENV CGO_ENABLED=0 GOOS=linux

# Compila o aplicativo
RUN go build -o /app/bin/app .

# Estágio final (imagem leve)
FROM alpine:latest AS final

# Define o diretório de trabalho
WORKDIR /root/

# Copia o binário compilado do estágio de construção
COPY --from=builder /app/bin/app .

# Copia o arquivo env-exemple para .env
COPY --from=builder /app/env-exemple .env
RUN chmod 644 .env

# Expõe a porta 3009 (se necessário)
EXPOSE 3009

# Define o usuário que executará o aplicativo (não-root)
USER 1001

# Comando padrão para executar o aplicativo
CMD ["./app"]