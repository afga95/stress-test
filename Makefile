# Makefile para Sistema CLI de Teste de Carga

# Variáveis
BINARY_NAME=load-test-cli
DOCKER_IMAGE=load-test-cli
GO_VERSION=1.21

# Comandos padrão
.PHONY: all build clean run docker-build docker-run test help

# Comando padrão
all: build

# Compilar a aplicação
build:
	@echo "Compilando a aplicação..."
	go build -o $(BINARY_NAME) main.go
	@echo "Aplicação compilada com sucesso: $(BINARY_NAME)"

# Executar a aplicação localmente
run:
	@echo "Executando a aplicação..."
	./$(BINARY_NAME) --url=http://google.com --requests=100 --concurrency=10

# Limpar arquivos compilados
clean:
	@echo "Limpando arquivos compilados..."
	rm -f $(BINARY_NAME)
	@echo "Limpeza concluída."

# Construir imagem Docker
docker-build:
	@echo "Construindo imagem Docker..."
	docker build -t $(DOCKER_IMAGE) .
	@echo "Imagem Docker construída: $(DOCKER_IMAGE)"

# Executar com Docker (exemplo básico)
docker-run:
	@echo "Executando com Docker..."
	docker run $(DOCKER_IMAGE) --url=http://google.com --requests=1000 --concurrency=10

# Executar teste customizado com Docker
docker-test:
	@echo "Executando teste customizado..."
	docker run $(DOCKER_IMAGE) --url=https://httpbin.org/get --requests=200 --concurrency=15

# Executar teste de alta concorrência
docker-stress:
	@echo "Executando teste de stress..."
	docker run $(DOCKER_IMAGE) --url=https://httpbin.org/delay/1 --requests=1000 --concurrency=50

# Instalar dependências Go
deps:
	@echo "Instalando dependências..."
	go mod init load-test-cli 2>/dev/null || true
	go mod tidy
	@echo "Dependências instaladas."

# Formatação do código
fmt:
	@echo "Formatando código..."
	go fmt ./...

# Verificação de lint
lint:
	@echo "Verificando lint..."
	golint ./...

# Executar testes unitários (se houver)
test:
	@echo "Executando testes..."
	go test ./...

# Mostrar informações da versão
version:
	@echo "Go version: $(shell go version)"
	@echo "Binary: $(BINARY_NAME)"
	@echo "Docker image: $(DOCKER_IMAGE)"

# Remover imagem Docker
docker-clean:
	@echo "Removendo imagem Docker..."
	docker rmi $(DOCKER_IMAGE) 2>/dev/null || true
	@echo "Imagem removida."

# Exemplo completo de uso
example:
	@echo "Exemplos de uso:"
	@echo ""
	@echo "1. Compilar e executar localmente:"
	@echo "   make build && ./$(BINARY_NAME) --url=http://google.com --requests=100 --concurrency=10"
	@echo ""
	@echo "2. Executar com Docker:"
	@echo "   make docker-build && make docker-run"
	@echo ""
	@echo "3. Teste customizado:"
	@echo "   docker run $(DOCKER_IMAGE) --url=https://api.exemplo.com --requests=500 --concurrency=25"

# Ajuda
help:
	@echo "Comandos disponíveis:"
	@echo "  build         - Compilar a aplicação"
	@echo "  run           - Executar a aplicação localmente"
	@echo "  clean         - Limpar arquivos compilados"
	@echo "  docker-build  - Construir imagem Docker"
	@echo "  docker-run    - Executar com Docker"
	@echo "  docker-test   - Executar teste customizado"
	@echo "  docker-stress - Executar teste de stress"
	@echo "  docker-clean  - Remover imagem Docker"
	@echo "  deps          - Instalar dependências"
	@echo "  fmt           - Formatar código"
	@echo "  lint          - Verificar lint"
	@echo "  test          - Executar testes"
	@echo "  version       - Mostrar versões"
	@echo "  example       - Mostrar exemplos"
	@echo "  help          - Mostrar esta ajuda"