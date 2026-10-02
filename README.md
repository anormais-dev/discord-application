# Roleta de Times

Discord Activity que monta times por sorteio numa roleta. Backend em Go (`backend/`), frontend em React + Rsbuild (`frontend/`).

## Mexer no frontend sem Go e sem Discord

Precisa só do Docker.

```bash
docker compose up
```

Abra http://localhost:8080, digite um nome e pronto. Para simular várias pessoas na mesma roleta, abra outras abas (ou janelas anônimas) com nomes diferentes. Também dá para entrar direto pela URL: `http://localhost:8080/?user=Ana`.

Edite os arquivos em `frontend/src` normalmente. O Rsbuild recompila (já minificado) em poucos segundos; depois é só recarregar a página.

- Checar tipos: `docker compose exec frontend npm run typecheck`
- Porta 8080 ocupada: `APP_PORT=8081 docker compose up`
- Dependência nova no `package.json`: `docker compose restart frontend` (ele roda `npm ci` ao subir)

Nesse modo backend e frontend rodam com `DEV_AUTH=true`, que aceita qualquer nome sem login. Sem essa variável, quem abrir o link fora do Discord vê só um aviso para entrar pela Activity. Nunca ative isso em produção.

## Rodar dentro do Discord

1. Defina as variáveis de ambiente `PORT` (obrigatória), `DISCORD_CLIENT_ID` e `DISCORD_CLIENT_SECRET` (as duas últimas vêm do Developer Portal, aba **OAuth2**; não é o token do bot). Backend e frontend leem só do ambiente; o `make` exporta o que estiver num `.env` na raiz, então dá para copiar o `.env.example` para `.env` e preencher.
2. No Developer Portal: ative **Activities**, adicione o redirect `https://127.0.0.1` em OAuth2 e, em **URL Mappings**, aponte `/` para o domínio do tunnel.
3. `make dev` (precisa de Go e Node).
4. `cloudflared tunnel --protocol http2 --url http://localhost:<PORT>` e use o domínio gerado no URL Mapping.
5. Abra a Activity num canal de voz, pelo ícone de foguete.

## Variáveis de ambiente

Ficam no `.env` (copie de `.env.example`).

| Variável | Obrigatória | Para que serve |
|---|---|---|
| `DISCORD_CLIENT_ID` | Para rodar no Discord | ID da aplicação, na aba **OAuth2** do Developer Portal. Também vai embutido no bundle do frontend durante o build. |
| `DISCORD_CLIENT_SECRET` | Para rodar no Discord | Secret da aba **OAuth2**. Não é o token do bot. |
| `PORT` | Sim | Porta HTTP do backend. Padrão do exemplo: `3000`. |
| `DEV_AUTH` | Não | Com `true`, aceita qualquer nome sem login do Discord. Só para desenvolvimento, nunca em produção. |

Duas variáveis opcionais não aparecem no `.env.example` porque o padrão já funciona:

- `STATIC_DIR`: pasta com o build do frontend que o backend serve. Padrão: `../frontend/dist`.
- `APP_PORT`: porta exposta pelo `docker compose up`. Padrão: `8080`.

## Testes

```bash
make test
```
