# Roleta de Times

Discord Activity que monta times por sorteio numa roleta. Backend em Go (`backend/`), frontend em React + Rsbuild (`frontend/`).

## Mexer no frontend sem Go e sem Discord

Precisa só do Docker. O compose lê `PORT` e `DEV_AUTH` do `.env`, então copie o `.env.example` e preencha `PORT` e `DEV_AUTH=true`.

```bash
docker compose up -d
```

Abra http://localhost:8080, digite um nome e pronto. Para simular várias pessoas na mesma roleta, abra outras abas (ou janelas anônimas) com nomes diferentes. Também dá para entrar direto pela URL: `http://localhost:8080/?user=Ana`.

Edite os arquivos em `frontend/src` normalmente. O Rsbuild recompila (já minificado) em poucos segundos; depois é só recarregar a página.

- Ver os logs: `docker compose logs -f`
- Parar: `docker compose down`
- Checar tipos: `docker compose exec frontend npm run typecheck`
- Porta 8080 ocupada: `APP_PORT=8081 docker compose up -d`
- Dependência nova no `package.json`: `docker compose restart frontend` (ele roda `npm ci` ao subir)

Com `DEV_AUTH=true`, backend e frontend aceitam qualquer nome sem login. Sem essa variável, quem abrir o link fora do Discord vê só um aviso para entrar pela Activity. Nunca ative isso em produção.

## Rodar dentro do Discord

1. Defina as variáveis de ambiente `PORT` (obrigatória), `DISCORD_CLIENT_ID` e `DISCORD_CLIENT_SECRET` (as duas últimas vêm do Developer Portal, aba **OAuth2**; não é o token do bot). O backend lê do ambiente e o frontend busca o client ID no backend (`GET /api/config`); o `make` exporta o que estiver num `.env` na raiz, então dá para copiar o `.env.example` para `.env` e preencher.
2. No Developer Portal: ative **Activities**, adicione o redirect `https://127.0.0.1` em OAuth2 e, em **URL Mappings**, aponte `/` para o domínio do tunnel.
3. `make dev` (precisa de Go e Node).
4. `cloudflared tunnel --protocol http2 --url http://localhost:<PORT>` e use o domínio gerado no URL Mapping.
5. Abra a Activity num canal de voz, pelo ícone de foguete.

## Variáveis de ambiente

Ficam no `.env` (copie de `.env.example`).

| Variável | Obrigatória | Para que serve |
|---|---|---|
| `DISCORD_CLIENT_ID` | Para rodar no Discord | ID da aplicação, na aba **OAuth2** do Developer Portal. Só o backend lê; o frontend recebe pelo `GET /api/config`, então o build não precisa dele. |
| `DISCORD_CLIENT_SECRET` | Para rodar no Discord | Secret da aba **OAuth2**. Não é o token do bot. |
| `PORT` | Sim | Porta HTTP do backend. Não tem padrão no código nem na imagem; o exemplo usa `3000`. |
| `DEV_AUTH` | Não | Com `true`, aceita qualquer nome sem login do Discord. Só para desenvolvimento, nunca em produção. |

Duas variáveis opcionais não aparecem no `.env.example` porque o padrão já funciona:

- `STATIC_DIR`: pasta com o build do frontend que o backend serve. Padrão: `../frontend/dist`.
- `APP_PORT`: porta exposta pelo `docker compose up -d`. Padrão: `8080`.

## Produção

Os arquivos de produção ficam em `infra/`. O `docker-compose.yml` e o `backend/Dockerfile` são só de desenvolvimento.

| Arquivo | Para que serve |
|---|---|
| `infra/backend.Dockerfile` | Imagem do backend: só a API (`/api/config`, `/api/token` e `/api/ws`). |
| `infra/frontend.Dockerfile` | Imagem do frontend: nginx servindo o build do Rsbuild. |
| `infra/nginx.conf` | Cache (`index.html` sem cache, `/static/` imutável), gzip e headers básicos. |
| `infra/stack.yml` | Stack que o Dokploy roda (Compose no modo Docker Stack). |

As imagens usam a raiz do repo como contexto. Cada Dockerfile tem um `.dockerignore` próprio ao lado, que manda para o build só a pasta daquele app.

```bash
docker build -f infra/backend.Dockerfile -t backend .
docker build -f infra/frontend.Dockerfile -t frontend .
```

Os builds são multi-stage, em três estágios:

1. `deps`: instala as dependências (`go mod download` ou `npm ci`). Fica em cache enquanto `go.mod`/`package-lock.json` não mudarem.
2. `build`: compila o binário Go ou o bundle do frontend. Roda na arquitetura de quem builda e gera para a do alvo, então dá para gerar `linux/arm64` sem emulação.
3. `runtime`: imagem final mínima, só com o resultado do build e sem root.

No servidor, o Traefik do Dokploy recebe as requisições e separa por caminho: `/api` vai para o backend e o resto para o frontend. O `stack.yml` lê do ambiente do serviço no Dokploy:

| Variável | Para que serve |
|---|---|
| `DOMAIN` | Domínio público do app, usado nas regras do Traefik. |
| `PORT` | Porta do backend dentro do container. |
| `DISCORD_CLIENT_ID` | Igual ao do desenvolvimento. |
| `DISCORD_CLIENT_SECRET` | Igual ao do desenvolvimento. |
| `IMAGE_TAG` | Opcional. Tag das imagens no GHCR; padrão `latest`. |

`DEV_AUTH` não existe em produção.

### Deploy

O deploy é automático a cada push na `main`:

1. O workflow `build` builda só os apps que mudaram e publica as imagens no GHCR com as tags `latest` e o SHA curto do commit.
2. O job `deploy` chama o webhook do Dokploy, que redeploya o stack com o `latest` novo.

Se nenhum app mudou, nada é buildado e o deploy não roda.

Para funcionar:

- o secret `DEPLOY_WEBHOOK` do repositório guarda a Webhook URL do serviço (aba Deployments no Dokploy);
- o Autodeploy do serviço fica ligado, senão o Dokploy recusa o webhook;
- o `IMAGE_TAG` fica vazio ou como `latest`.

Para fazer rollback, coloque `IMAGE_TAG=<sha curto>` de um commit anterior no Dokploy e clique em Deploy. Enquanto o SHA estiver fixado, os próximos pushes redeployam ele mesmo. Depois da correção, volte para `latest`.

## Testes

```bash
make test
```
