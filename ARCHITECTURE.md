# Arquitetura do Sistema de Recrutamento e Seleção

## Visão Geral

O sistema segue a arquitetura hexagonal (Ports & Adapters) e é dividido em microserviços independentes que se comunicam via HTTP REST APIs.

## Diagrama da Arquitetura

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              FRONTEND / CLIENT                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      │ HTTP/REST
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              API GATEWAY                                    │
│                         (Load Balancer/Nginx)                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                    ┌─────────────────┼─────────────────┐
                    │                 │                 │
                    ▼                 ▼                 ▼
        ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
        │  AUTH SERVICE   │ │  JOB SERVICE    │ │CANDIDATE SERVICE│
        │    Port 8083    │ │   Port 8081     │ │   Port 8082     │
        └─────────────────┘ └─────────────────┘ └─────────────────┘
                    │                 │                 │
                    └─────────────────┼─────────────────┘
                                      │
                                      ▼
                    ┌─────────────────────────────────────┐
                    │          PostgreSQL Database        │
                    │            Port 5432                │
                    └─────────────────────────────────────┘
```

## Microserviços

### 1. Auth Service (Port 8083)
**Responsabilidades:**
- Autenticação e autorização de usuários
- Gerenciamento de tokens JWT
- Registro de usuários (admin/candidatos)
- Validação de tokens para outros serviços

**Endpoints:**
- `POST /api/v1/auth/register` - Registrar usuário
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/auth/profile` - Perfil do usuário
- `PUT /api/v1/auth/change-password` - Alterar senha
- `POST /api/v1/auth/validate` - Validar token

### 2. Job Service (Port 8081)
**Responsabilidades:**
- CRUD de vagas de trabalho
- Gerenciamento de status das vagas
- Associação de skills às vagas
- Busca e filtros de vagas

**Endpoints:**
- `POST /api/v1/jobs` - Criar vaga
- `GET /api/v1/jobs` - Listar vagas
- `GET /api/v1/jobs/:id` - Obter vaga por ID
- `PUT /api/v1/jobs/:id` - Atualizar vaga
- `DELETE /api/v1/jobs/:id` - Excluir vaga
- `PATCH /api/v1/jobs/:id/status` - Alterar status

### 3. Candidate Service (Port 8082)
**Responsabilidades:**
- Gerenciamento de perfis de candidatos
- Upload e processamento de currículos
- Candidaturas às vagas
- Integração com IA para análise de currículos

**Endpoints:**
- `POST /api/v1/candidates` - Registrar candidato
- `GET /api/v1/candidates/:id` - Obter candidato
- `PUT /api/v1/candidates/:id` - Atualizar candidato
- `POST /api/v1/candidates/:id/resume` - Upload currículo
- `POST /api/v1/candidates/:id/applications` - Candidatar-se

## Arquitetura Hexagonal por Serviço

Cada microserviço segue a estrutura:

```
services/[service-name]/
├── cmd/
│   └── main.go                 # Entry point
├── internal/
│   ├── domain/                 # Entidades e regras de negócio
│   │   ├── entities.go
│   │   └── repository.go       # Ports (interfaces)
│   ├── application/            # Casos de uso
│   │   └── service.go
│   ├── infrastructure/         # Adapters (implementações)
│   │   ├── repository.go
│   │   └── external_apis.go
│   └── interfaces/             # Controllers e rotas
│       ├── controller.go
│       └── routes.go
└── Dockerfile
```

## Camadas da Arquitetura Hexagonal

### Domain (Centro)
- **Entidades**: Modelos de negócio puros
- **Value Objects**: Objetos imutáveis
- **Ports**: Interfaces que definem contratos
- **Regras de Negócio**: Lógica core da aplicação

### Application
- **Use Cases**: Orquestração das regras de negócio
- **Services**: Serviços de aplicação
- **DTOs**: Objetos de transferência de dados

### Infrastructure (Adapters)
- **Repositories**: Implementação de persistência
- **External APIs**: Integração com serviços externos
- **Database**: Configuração e conexão com BD

### Interfaces (Adapters)
- **Controllers**: Handlers HTTP
- **Routes**: Definição de rotas
- **Middleware**: Interceptadores de requisições

## Banco de Dados

### Tabelas Principais

```sql
users                    # Usuários (admin/candidatos)
├── id (UUID)
├── email
├── password_hash
├── role
└── name

jobs                     # Vagas de trabalho
├── id (UUID)
├── title
├── description
├── requirements
├── status
└── created_by (FK users)

candidates               # Perfis de candidatos
├── id (UUID)
├── user_id (FK users)
├── phone
├── address
└── linkedin_url

job_applications         # Candidaturas
├── id (UUID)
├── job_id (FK jobs)
├── candidate_id (FK candidates)
├── status
└── cover_letter

skills                   # Habilidades
├── id (UUID)
├── name
└── category

candidate_skills         # Skills dos candidatos
├── candidate_id (FK)
├── skill_id (FK)
├── proficiency_level
└── years_of_experience
```

## Comunicação Entre Serviços

### Padrões de Comunicação
1. **HTTP REST**: Comunicação síncrona entre serviços
2. **Database Sharing**: Compartilhamento de dados via BD
3. **Event-Driven** (futuro): Mensageria assíncrona

### Fluxos de Comunicação

#### Autenticação
```
Client → Auth Service → JWT Token
Client → Other Services (with JWT) → Auth Service (validate)
```

#### Candidatura a Vaga
```
Client → Candidate Service → Job Service (check job status)
                          → Auth Service (validate user)
```

## Segurança

### Autenticação
- JWT tokens com expiração
- Refresh tokens para renovação
- Middleware de autenticação compartilhado

### Autorização
- Role-based access control (RBAC)
- Admin: Gerenciar vagas
- Candidate: Gerenciar perfil e candidaturas

### Validação
- Validação de entrada em todos os endpoints
- Sanitização de dados
- Rate limiting (futuro)

## Deployment

### Docker
- Cada serviço tem seu próprio Dockerfile
- Docker Compose para orquestração local
- Multi-stage builds para otimização

### Ambiente de Produção
- Kubernetes (recomendado)
- Load balancer (Nginx/HAProxy)
- Database clustering
- Monitoring e logging

## Escalabilidade

### Horizontal Scaling
- Cada serviço pode ser escalado independentemente
- Load balancing entre instâncias
- Database connection pooling

### Performance
- Índices otimizados no banco
- Cache (Redis) para dados frequentes
- CDN para arquivos estáticos

## Monitoramento

### Métricas
- Health checks em cada serviço
- Métricas de performance
- Logs estruturados

### Observabilidade
- Distributed tracing
- Error tracking
- Performance monitoring
