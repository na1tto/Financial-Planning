# Go Financial Planning

Aplicação web para planejamento financeiro com:
- cadastro e login de usuários
- registro de receitas/despesas futuras
- gráficos de projeção mensal
- autenticação stateful via sessão em cookie

## Tecnologias

- Backend: Go + Chi + SQLite + Goose
- Frontend: React + Vite + Recharts
- Segurança: Argon2id (senhas), cookie `HttpOnly`, sessão persistida no banco

## Pré-requisitos

Antes de rodar, tenha instalado:

- Go 1.24+ (ou compatível com o `go.mod`)
- Node.js 20+ e npm

Verifique:

```bash
go version
node -v
npm -v
```

## Estrutura do projeto

```text
backend/
  cmd/api/main.go
  internal/
    app/
    auth/
    db/
    handlers/
    repository/
frontend/
  src/
```

## Como rodar (passo a passo)

Abra **2 terminais** na raiz do projeto.

### 1) Backend (terminal 1)

```bash
cd backend
go mod tidy
go run ./cmd/api/main.go
```

Ao subir corretamente, deve aparecer algo como:

```text
goose: no migrations to run. current version: 1
api listening on :8080
```

Observações:
- As migrations rodam automaticamente no startup.
- O banco SQLite é criado em `backend/finance.db` (por padrão).
- O projeto usa driver SQLite em Go puro (`modernc.org/sqlite`), então não exige CGO.

### 2) Frontend (terminal 2)

```bash
cd frontend
npm install
npm run dev
```

Depois acesse:

- `http://localhost:5173`

## Variáveis de ambiente

### Backend

- `API_ADDR` (default: `:8080`)
- `DB_PATH` (default: `./finance.db`)
- `FRONTEND_ORIGIN` (default: `http://localhost:5173`)
- `APP_PRODUCTION` (default: `false`)
- `SESSION_TTL` (default: `168h`)

Exemplo (PowerShell):

```powershell
$env:API_ADDR=":8090"
$env:FRONTEND_ORIGIN="http://localhost:5173"
go run ./cmd/api/main.go
```

### Frontend

- `VITE_API_URL` (default: `http://localhost:8080/api/v1`)

Se backend estiver em outra porta (exemplo `8090`), configure:

```powershell
$env:VITE_API_URL="http://localhost:8090/api/v1"
npm run dev
```

## Endpoints da API

### Healthcheck

- `GET /healthz`

### Autenticação

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/me`

### Transações

- `GET /api/v1/transactions`
- `POST /api/v1/transactions`
- `PUT /api/v1/transactions/{id}`
- `DELETE /api/v1/transactions/{id}`

### Projeção

- `GET /api/v1/forecast/monthly?months=12`

## Regras de payload

- `kind`: `income` ou `expense`
- `amount`: decimal positivo (ex.: `1299.90`)
- `dueDate`: `YYYY-MM-DD`
- Login aceita apenas `email` e `password`
- Cadastro aceita `name`, `email`, `password`

## Solução de problemas

### 1) Erro de porta em uso (`bind: ... 8080`)

Rode backend em outra porta:

```powershell
$env:API_ADDR=":8090"
go run ./cmd/api/main.go
```

E ajuste o frontend:

```powershell
$env:VITE_API_URL="http://localhost:8090/api/v1"
npm run dev
```

### 2) Erro de CORS

Garanta que `FRONTEND_ORIGIN` no backend aponta para a URL exata do frontend (incluindo porta).

### 3) Erro `unknown field: name` no login

Use login com payload apenas:

```json
{
  "email": "usuario@exemplo.com",
  "password": "senha1234"
}
```

### 4) Frontend aparentemente não atualiza

Pare e suba novamente:

```bash
npm run dev
```

Se necessário, limpe build anterior:

```bash
cd frontend
npm run build
```

## Segurança implementada

- Senhas com hash Argon2id
- Sessão stateful com token aleatório e hash SHA-256 no banco
- Cookie com `HttpOnly` e `SameSite=Lax`
- Headers de segurança (`X-Frame-Options`, `X-Content-Type-Options`, `CSP`)
