# Guia de Instalação e Execução

## Pré-requisitos

### Desenvolvimento Local
- Go 1.21 ou superior
- PostgreSQL 12 ou superior
- Docker e Docker Compose (opcional, mas recomendado)
- Git

### Apenas Docker
- Docker
- Docker Compose

## Instalação

### 1. Clone o Repositório

```bash
git clone <repository-url>
cd RecrutamentoSelecao-Backend-Go
```

### 2. Configuração do Ambiente

Copie o arquivo de exemplo e configure as variáveis:

```bash
cp .env.example .env
```

Edite o arquivo `.env` com suas configurações:

```bash
# Exemplo de configuração local
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=recruitment_db
JWT_SECRET=your-super-secret-jwt-key
```

## Execução com Docker (Recomendado)

### 1. Executar todos os serviços

```bash
# Construir e executar todos os serviços
docker-compose up -d

# Verificar logs
docker-compose logs -f

# Verificar status dos serviços
docker-compose ps
```

### 2. Executar migrations (se necessário)

```bash
# As migrations são executadas automaticamente na inicialização do PostgreSQL
# Mas você pode executar manualmente se necessário
docker-compose exec postgres psql -U postgres -d recruitment_db -f /docker-entrypoint-initdb.d/001_init.sql
```

### 3. Parar os serviços

```bash
docker-compose down
```

## Execução Local (Desenvolvimento)

### 1. Configurar PostgreSQL

Instale e configure o PostgreSQL:

```bash
# Ubuntu/Debian
sudo apt-get install postgresql postgresql-contrib

# macOS
brew install postgresql

# Windows
# Baixe e instale do site oficial: https://www.postgresql.org/download/windows/
```

Crie o banco de dados:

```bash
sudo -u postgres psql
CREATE DATABASE recruitment_db;
CREATE USER postgres WITH PASSWORD 'postgres';
GRANT ALL PRIVILEGES ON DATABASE recruitment_db TO postgres;
\q
```

### 2. Executar migrations

```bash
psql -U postgres -d recruitment_db -f migrations/001_init.sql
```

### 3. Instalar dependências Go

```bash
go mod tidy
```

### 4. Executar os serviços

#### Opção 1: Executar todos os serviços (se make estiver disponível)

```bash
# No Linux/macOS com make
make run-all
```

#### Opção 2: Executar serviços individualmente

```bash
# Terminal 1 - Auth Service
cd services/auth-service
go run cmd/main.go

# Terminal 2 - Job Service (quando implementado)
cd services/job-service
go run cmd/main.go

# Terminal 3 - Candidate Service (quando implementado)
cd services/candidate-service
go run cmd/main.go
```

## Verificação da Instalação

### 1. Health Checks

Verifique se os serviços estão funcionando:

```bash
# Auth Service
curl http://localhost:8083/health

# Job Service
curl http://localhost:8081/health

# Candidate Service
curl http://localhost:8082/health
```

### 2. Teste de Registro e Login

```bash
# Registrar um usuário admin
curl -X POST http://localhost:8083/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "password123",
    "name": "Admin User",
    "role": "admin"
  }'

# Fazer login
curl -X POST http://localhost:8083/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "password123"
  }'
```

## Testes

### Executar todos os testes

```bash
# Com make
make test

# Ou diretamente com go
go test -v ./...
```

### Executar testes com coverage

```bash
# Com make
make test-coverage

# Ou diretamente com go
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Desenvolvimento

### Estrutura de Desenvolvimento

```bash
# Formatar código
go fmt ./...

# Verificar problemas
go vet ./...

# Executar linter (se golangci-lint estiver instalado)
golangci-lint run
```

### Hot Reload (Opcional)

Para desenvolvimento com hot reload, instale o Air:

```bash
go install github.com/cosmtrek/air@latest

# Executar com hot reload
air
```

## Troubleshooting

### Problemas Comuns

#### 1. Erro de conexão com banco de dados

```bash
# Verificar se PostgreSQL está rodando
sudo systemctl status postgresql

# Verificar conexão
psql -U postgres -d recruitment_db -c "SELECT 1;"
```

#### 2. Porta já em uso

```bash
# Verificar processos usando as portas
lsof -i :8083
lsof -i :8081
lsof -i :8082

# Matar processo se necessário
kill -9 <PID>
```

#### 3. Problemas com Docker

```bash
# Limpar containers e volumes
docker-compose down -v
docker system prune -f

# Reconstruir imagens
docker-compose build --no-cache
```

#### 4. Problemas com Go modules

```bash
# Limpar cache de módulos
go clean -modcache

# Reinstalar dependências
rm go.sum
go mod tidy
```

### Logs e Debug

#### Docker Logs

```bash
# Ver logs de todos os serviços
docker-compose logs

# Ver logs de um serviço específico
docker-compose logs auth-service
docker-compose logs postgres

# Seguir logs em tempo real
docker-compose logs -f auth-service
```

#### Logs Locais

Os serviços logam no stdout/stderr. Para debug mais detalhado, configure a variável de ambiente:

```bash
export LOG_LEVEL=debug
```

## Configurações Avançadas

### Configuração de Produção

Para ambiente de produção, considere:

1. **Segurança**:
   - Use senhas fortes
   - Configure HTTPS
   - Use secrets management
   - Configure rate limiting

2. **Performance**:
   - Configure connection pooling
   - Use cache (Redis)
   - Configure load balancer

3. **Monitoramento**:
   - Configure logging centralizado
   - Use métricas (Prometheus)
   - Configure alertas

### Configuração de Banco de Dados

Para melhor performance em produção:

```sql
-- Configurações recomendadas para PostgreSQL
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET default_statistics_target = 100;
```

## Suporte

Para problemas ou dúvidas:

1. Verifique os logs dos serviços
2. Consulte a documentação da API
3. Verifique issues conhecidos no repositório
4. Crie uma issue detalhando o problema
