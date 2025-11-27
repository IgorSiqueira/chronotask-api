# 🐳 ChronoTask API - Docker & Coolify Deploy

## 📋 Pré-requisitos

- Docker 24.0+
- Docker Compose 2.0+
- Coolify (para deploy em produção)

## 🚀 Deploy Local

### 1. Configurar Variáveis de Ambiente

```bash
# Copiar arquivo de exemplo
cp .env.example .env

# Editar .env com suas configurações
nano .env
```

**IMPORTANTE**: Altere os seguintes valores em produção:
- `DB_PASSWORD`: senha forte para PostgreSQL
- `JWT_SECRET`: chave secreta de 32+ caracteres (use: `openssl rand -base64 32`)

### 2. Subir os Containers

```bash
# Build e start
docker-compose up -d

# Ver logs
docker-compose logs -f api

# Verificar status
docker-compose ps
```

### 3. Verificar Health

```bash
# Health check da API
curl http://localhost:8080/health

# Resposta esperada:
# {"status":"ok"}
```

### 4. Testar Endpoints

```bash
# Criar usuário
curl -X POST http://localhost:8080/api/v1/user \
  -H "Content-Type: application/json" \
  -d '{
    "fullName": "João Silva",
    "email": "joao@example.com",
    "password": "senha123",
    "birthDate": "1990-01-15",
    "acceptTerms": true
  }'

# Login
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "joao@example.com",
    "password": "senha123"
  }'
```

## ☁️ Deploy no Coolify

### 1. Configurar Projeto no Coolify

1. Acesse seu Coolify
2. Criar novo Resource → Docker Compose
3. Cole o conteúdo do `docker-compose.yml`

### 2. Configurar Variáveis de Ambiente

No Coolify, adicione as seguintes variáveis:

```env
DB_PASSWORD=<senha_forte_postgres>
JWT_SECRET=<chave_jwt_forte>
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=168h
PORT=8080
GIN_MODE=release
```

### 3. Configurar Domínio

1. Coolify → Settings → Domains
2. Adicionar domínio: `api.chronotask.com`
3. Habilitar HTTPS automático
4. Coolify gerará certificado SSL via Let's Encrypt

### 4. Deploy

1. Clicar em "Deploy"
2. Coolify irá:
   - Fazer pull do repositório
   - Build da imagem Docker
   - Subir containers
   - Configurar proxy reverso
   - Gerar certificado SSL

### 5. Verificar Deploy

```bash
# Health check
curl https://api.chronotask.com/health

# Resposta: {"status":"ok"}
```

## 🔧 Comandos Úteis

### Gerenciar Containers

```bash
# Parar containers
docker-compose down

# Parar e remover volumes (CUIDADO: apaga dados)
docker-compose down -v

# Rebuild sem cache
docker-compose build --no-cache

# Ver logs de um serviço específico
docker-compose logs -f postgres
docker-compose logs -f api

# Executar comando no container
docker-compose exec api /bin/sh
docker-compose exec postgres psql -U postgres -d chronotask
```

### Database

```bash
# Acessar PostgreSQL
docker-compose exec postgres psql -U postgres -d chronotask

# Backup do banco
docker-compose exec postgres pg_dump -U postgres chronotask > backup.sql

# Restore do banco
docker-compose exec -T postgres psql -U postgres chronotask < backup.sql

# Ver migrations aplicadas
docker-compose exec postgres psql -U postgres -d chronotask -c "SELECT * FROM schema_migrations;"
```

### Monitoramento

```bash
# Status dos containers
docker-compose ps

# Uso de recursos
docker stats

# Inspecionar health
docker inspect --format='{{json .State.Health}}' chronotask-api | jq
```

## 📊 Arquitetura dos Containers

```
┌─────────────────────────────────────────┐
│         Coolify Proxy (Traefik)        │
│         HTTPS/SSL Automático            │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│        ChronoTask API (Go)              │
│        Port: 8080                       │
│        Health: /health                  │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│        PostgreSQL 16                    │
│        Port: 5432                       │
│        Volume: postgres_data            │
└─────────────────────────────────────────┘
```

## 🔒 Segurança

### Checklist de Produção

- [ ] `DB_PASSWORD` forte (16+ caracteres, letras, números, símbolos)
- [ ] `JWT_SECRET` único e forte (32+ caracteres)
- [ ] `GIN_MODE=release` em produção
- [ ] HTTPS habilitado (Coolify faz automaticamente)
- [ ] Firewall configurado (apenas portas 80/443 expostas)
- [ ] Backup automático do PostgreSQL
- [ ] Logs centralizados
- [ ] Monitoramento ativo (uptime, performance)

### Gerar Senhas Fortes

```bash
# JWT Secret
openssl rand -base64 48

# Database Password
openssl rand -base64 32 | tr -d "=+/" | cut -c1-25
```

## 🐛 Troubleshooting

### API não inicia

```bash
# Ver logs detalhados
docker-compose logs api

# Verificar se PostgreSQL está healthy
docker-compose ps postgres

# Testar conexão com banco
docker-compose exec api wget -O- http://postgres:5432
```

### Erro de conexão com banco

```bash
# Verificar variáveis de ambiente
docker-compose exec api env | grep DB_

# Testar conexão manualmente
docker-compose exec api psql -h postgres -U postgres -d chronotask
```

### Migrations não aplicadas

```bash
# Verificar se migrations estão no container
docker-compose exec api ls -la /app/migrations

# Aplicar migrations manualmente
docker-compose exec postgres psql -U postgres -d chronotask -f /docker-entrypoint-initdb.d/001_create_users_table.sql
```

## 📈 Performance

### Otimizações de Produção

1. **Connection Pooling**: Configurado no código (25 max connections)
2. **Multi-stage Build**: Imagem final ~20MB
3. **Health Checks**: Monitoramento automático
4. **Non-root User**: Segurança aprimorada
5. **Static Binary**: Sem dependências runtime

### Recursos Recomendados

**Produção**:
- CPU: 1 vCPU
- RAM: 512MB (API) + 512MB (PostgreSQL)
- Disk: 10GB (para logs e dados)

**Alto Tráfego**:
- CPU: 2 vCPU
- RAM: 1GB (API) + 2GB (PostgreSQL)
- Disk: 50GB SSD

## 📚 Referências

- [Docker Compose Docs](https://docs.docker.com/compose/)
- [Coolify Docs](https://coolify.io/docs)
- [PostgreSQL Docker](https://hub.docker.com/_/postgres)
- [Go Docker Best Practices](https://docs.docker.com/language/golang/)
