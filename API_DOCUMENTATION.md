# Documentação da API

## Visão Geral

O sistema de recrutamento e seleção expõe APIs REST para gerenciar usuários, vagas e candidatos. Todas as APIs seguem os padrões REST e retornam dados em formato JSON.

## Base URLs

- **Auth Service**: `http://localhost:8083/api/v1`
- **Job Service**: `http://localhost:8081/api/v1`
- **Candidate Service**: `http://localhost:8082/api/v1`

## Autenticação

A maioria dos endpoints requer autenticação via JWT token no header:

```
Authorization: Bearer <jwt_token>
```

## Formato de Resposta

Todas as respostas seguem o formato padrão:

```json
{
  "success": true,
  "message": "Operação realizada com sucesso",
  "data": { ... },
  "error": null
}
```

Para respostas paginadas:

```json
{
  "success": true,
  "message": "Dados recuperados com sucesso",
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

## Auth Service API

### Registrar Usuário

**POST** `/auth/register`

Registra um novo usuário no sistema.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "Nome do Usuário",
  "role": "admin" // ou "candidate"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "Nome do Usuário",
    "role": "admin"
  }
}
```

**Status Codes:**
- `201`: Usuário criado com sucesso
- `400`: Dados inválidos ou usuário já existe

### Login

**POST** `/auth/login`

Autentica um usuário e retorna tokens de acesso.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "jwt_access_token",
    "refresh_token": "refresh_token",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "name": "Nome do Usuário",
      "role": "admin"
    },
    "expires_at": "2024-01-01T12:00:00Z"
  }
}
```

**Status Codes:**
- `200`: Login realizado com sucesso
- `401`: Credenciais inválidas

### Refresh Token

**POST** `/auth/refresh`

Renova o token de acesso usando o refresh token.

**Request Body:**
```json
{
  "refresh_token": "refresh_token_here"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "token": "new_jwt_access_token",
    "refresh_token": "new_refresh_token",
    "user": { ... },
    "expires_at": "2024-01-01T12:00:00Z"
  }
}
```

### Validar Token

**POST** `/auth/validate`

Valida um token JWT.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Token is valid",
  "data": {
    "user_id": "uuid",
    "email": "user@example.com",
    "role": "admin"
  }
}
```

### Logout

**POST** `/auth/logout`

Invalida os tokens do usuário.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Logout successful"
}
```

### Obter Perfil

**GET** `/auth/profile`

Retorna informações do usuário autenticado.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "Nome do Usuário",
    "role": "admin"
  }
}
```

### Alterar Senha

**PUT** `/auth/change-password`

Altera a senha do usuário autenticado.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "current_password": "senha_atual",
  "new_password": "nova_senha"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Password changed successfully"
}
```

## Job Service API

### Criar Vaga

**POST** `/jobs`

Cria uma nova vaga de trabalho. Requer role `admin`.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "title": "Desenvolvedor Go",
  "description": "Descrição da vaga...",
  "requirements": "Requisitos da vaga...",
  "location": "São Paulo, SP",
  "salary_min": 5000.00,
  "salary_max": 8000.00,
  "skills": [
    {
      "skill_id": "uuid",
      "required_level": "intermediate",
      "is_required": true
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Job created successfully",
  "data": {
    "id": "uuid",
    "title": "Desenvolvedor Go",
    "description": "Descrição da vaga...",
    "status": "open",
    "created_by": "uuid",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

### Listar Vagas

**GET** `/jobs`

Lista vagas com paginação e filtros opcionais.

**Query Parameters:**
- `page`: Página (padrão: 1)
- `limit`: Itens por página (padrão: 10, máximo: 100)
- `status`: Filtrar por status (`open`, `closed`)
- `location`: Filtrar por localização
- `title`: Buscar por título

**Response:**
```json
{
  "success": true,
  "message": "Jobs retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "title": "Desenvolvedor Go",
      "description": "Descrição...",
      "location": "São Paulo, SP",
      "salary_min": 5000.00,
      "salary_max": 8000.00,
      "status": "open",
      "created_at": "2024-01-01T12:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 50,
    "total_pages": 5
  }
}
```

### Obter Vaga por ID

**GET** `/jobs/{id}`

Retorna detalhes de uma vaga específica.

**Response:**
```json
{
  "success": true,
  "message": "Job retrieved successfully",
  "data": {
    "id": "uuid",
    "title": "Desenvolvedor Go",
    "description": "Descrição completa...",
    "requirements": "Requisitos...",
    "location": "São Paulo, SP",
    "salary_min": 5000.00,
    "salary_max": 8000.00,
    "status": "open",
    "skills": [
      {
        "skill": {
          "id": "uuid",
          "name": "Go",
          "category": "Programming Language"
        },
        "required_level": "intermediate",
        "is_required": true
      }
    ],
    "created_by": "uuid",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

### Atualizar Vaga

**PUT** `/jobs/{id}`

Atualiza uma vaga existente. Requer role `admin` e ser o criador da vaga.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "title": "Desenvolvedor Go Senior",
  "description": "Nova descrição...",
  "requirements": "Novos requisitos...",
  "location": "São Paulo, SP",
  "salary_min": 6000.00,
  "salary_max": 10000.00
}
```

### Alterar Status da Vaga

**PATCH** `/jobs/{id}/status`

Altera o status de uma vaga. Requer role `admin`.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "status": "closed"
}
```

### Excluir Vaga

**DELETE** `/jobs/{id}`

Exclui uma vaga. Requer role `admin` e ser o criador da vaga.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

## Candidate Service API

### Registrar Candidato

**POST** `/candidates`

Registra um novo candidato no sistema.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "phone": "+55 11 99999-9999",
  "address": "Endereço completo",
  "date_of_birth": "1990-01-01",
  "linkedin_url": "https://linkedin.com/in/usuario",
  "github_url": "https://github.com/usuario"
}
```

### Obter Candidato

**GET** `/candidates/{id}`

Retorna informações de um candidato.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Candidate retrieved successfully",
  "data": {
    "id": "uuid",
    "user": {
      "id": "uuid",
      "email": "candidate@example.com",
      "name": "Nome do Candidato"
    },
    "phone": "+55 11 99999-9999",
    "address": "Endereço completo",
    "date_of_birth": "1990-01-01",
    "linkedin_url": "https://linkedin.com/in/usuario",
    "github_url": "https://github.com/usuario",
    "skills": [
      {
        "skill": {
          "id": "uuid",
          "name": "Go",
          "category": "Programming Language"
        },
        "proficiency_level": "advanced",
        "years_of_experience": 3
      }
    ],
    "work_experiences": [...],
    "education": [...],
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

### Upload de Currículo

**POST** `/candidates/{id}/resume`

Faz upload de um arquivo de currículo.

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: multipart/form-data
```

**Form Data:**
- `file`: Arquivo do currículo (PDF, DOC, DOCX)

**Response:**
```json
{
  "success": true,
  "message": "Resume uploaded successfully",
  "data": {
    "id": "uuid",
    "filename": "curriculo.pdf",
    "file_size": 1024000,
    "ai_processed": false,
    "uploaded_at": "2024-01-01T12:00:00Z"
  }
}
```

### Candidatar-se a Vaga

**POST** `/candidates/{id}/applications`

Candidata-se a uma vaga específica.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "job_id": "uuid",
  "cover_letter": "Carta de apresentação..."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Application submitted successfully",
  "data": {
    "id": "uuid",
    "job_id": "uuid",
    "candidate_id": "uuid",
    "status": "applied",
    "cover_letter": "Carta de apresentação...",
    "applied_at": "2024-01-01T12:00:00Z"
  }
}
```

## Skills API

### Listar Skills

**GET** `/skills`

Lista todas as skills disponíveis.

**Query Parameters:**
- `category`: Filtrar por categoria
- `search`: Buscar por nome

**Response:**
```json
{
  "success": true,
  "message": "Skills retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "name": "Go",
      "category": "Programming Language",
      "created_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

## Códigos de Status HTTP

- `200`: OK - Requisição bem-sucedida
- `201`: Created - Recurso criado com sucesso
- `400`: Bad Request - Dados inválidos
- `401`: Unauthorized - Token inválido ou ausente
- `403`: Forbidden - Permissões insuficientes
- `404`: Not Found - Recurso não encontrado
- `409`: Conflict - Conflito (ex: email já existe)
- `422`: Unprocessable Entity - Erro de validação
- `500`: Internal Server Error - Erro interno do servidor

## Exemplos de Uso

### Fluxo Completo de Candidatura

```bash
# 1. Registrar usuário candidato
curl -X POST http://localhost:8083/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "candidate@example.com",
    "password": "password123",
    "name": "João Silva",
    "role": "candidate"
  }'

# 2. Fazer login
TOKEN=$(curl -X POST http://localhost:8083/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "candidate@example.com",
    "password": "password123"
  }' | jq -r '.data.token')

# 3. Completar perfil de candidato
curl -X POST http://localhost:8082/api/v1/candidates \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "+55 11 99999-9999",
    "address": "São Paulo, SP",
    "linkedin_url": "https://linkedin.com/in/joaosilva"
  }'

# 4. Listar vagas disponíveis
curl -X GET http://localhost:8081/api/v1/jobs \
  -H "Authorization: Bearer $TOKEN"

# 5. Candidatar-se a uma vaga
curl -X POST http://localhost:8082/api/v1/candidates/CANDIDATE_ID/applications \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "job_id": "JOB_ID",
    "cover_letter": "Tenho interesse nesta vaga..."
  }'
```

## Rate Limiting

Atualmente não há rate limiting implementado, mas é recomendado para produção:

- Autenticação: 5 tentativas por minuto por IP
- APIs gerais: 100 requisições por minuto por usuário
- Upload de arquivos: 10 uploads por hora por usuário

## Versionamento

A API usa versionamento via URL path (`/api/v1/`). Mudanças breaking serão introduzidas em novas versões.

## Suporte

Para dúvidas sobre a API:
1. Consulte esta documentação
2. Verifique os exemplos de código
3. Teste os endpoints com as ferramentas fornecidas
4. Reporte bugs ou solicite melhorias via issues do repositório
