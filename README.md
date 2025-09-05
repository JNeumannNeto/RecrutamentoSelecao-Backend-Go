# Sistema de Recrutamento e Seleção - Backend Go

Sistema de gerenciamento de vagas de trabalho e candidatos desenvolvido em Go com arquitetura hexagonal e microserviços.

## Arquitetura

O sistema é composto por 3 microserviços principais:

1. **Job Service** - Gerenciamento de vagas
2. **Candidate Service** - Gerenciamento de candidatos e currículos
3. **Auth Service** - Autenticação e autorização de usuários

### Arquitetura Hexagonal

Cada microserviço segue a arquitetura hexagonal (Ports & Adapters) com as seguintes camadas:

- **Domain**: Entidades de negócio, value objects e regras de domínio
- **Application**: Casos de uso e serviços de aplicação
- **Infrastructure**: Adaptadores para banco de dados, APIs externas, etc.
- **Interfaces**: Controllers REST, handlers, etc.

## Tecnologias Utilizadas

- **Linguagem**: Go 1.21+
- **Banco de Dados**: PostgreSQL
- **Framework Web**: Gin
- **ORM**: GORM
- **Containerização**: Docker & Docker Compose
- **Testes**: Testify
- **Migrations**: golang-migrate

## Estrutura do Projeto

```
├── services/
│   ├── job-service/
│   ├── candidate-service/
│   └── auth-service/
├── shared/
│   ├── database/
│   ├── middleware/
│   └── utils/
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## Como Executar

### Pré-requisitos

- Docker e Docker Compose instalados
- Go 1.21+ (para desenvolvimento local)

### Executando com Docker

```bash
# Clone o repositório
git clone <repository-url>
cd RecrutamentoSelecao-Backend-Go

# Execute os serviços
docker-compose up -d

# Verifique os logs
docker-compose logs -f
```

### Executando Localmente

```bash
# Instale as dependências
go mod tidy

# Execute as migrations
make migrate-up

# Execute os testes
make test

# Execute os serviços
make run-all
```

## APIs

### Job Service (Port 8081)
- `POST /api/v1/jobs` - Criar vaga
- `GET /api/v1/jobs` - Listar vagas
- `GET /api/v1/jobs/:id` - Obter vaga por ID
- `PUT /api/v1/jobs/:id` - Atualizar vaga
- `DELETE /api/v1/jobs/:id` - Excluir vaga
- `PATCH /api/v1/jobs/:id/status` - Alterar status da vaga

### Candidate Service (Port 8082)
- `POST /api/v1/candidates` - Registrar candidato
- `GET /api/v1/candidates/:id` - Obter candidato por ID
- `PUT /api/v1/candidates/:id` - Atualizar candidato
- `POST /api/v1/candidates/:id/resume` - Upload de currículo
- `POST /api/v1/candidates/:id/applications` - Candidatar-se a vaga

### Auth Service (Port 8083)
- `POST /api/v1/auth/register` - Registrar usuário
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/logout` - Logout

## Banco de Dados

O sistema utiliza PostgreSQL com as seguintes tabelas principais:

- `users` - Usuários do sistema (admin/candidatos)
- `jobs` - Vagas de trabalho
- `candidates` - Informações dos candidatos
- `resumes` - Currículos dos candidatos
- `applications` - Candidaturas às vagas

## Testes

```bash
# Executar todos os testes
make test

# Executar testes com coverage
make test-coverage

# Executar testes de integração
make test-integration
```

## Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request
