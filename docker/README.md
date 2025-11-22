## Run multiple docker compose files

```bash
docker compose -f docker-compose.db.yml -f docker-compose.auth.yml -f docker-compose.tools.yml -f docker-compose.workers.yml up
```

## Run all docker compose files

```bash
docker compose $(ls docker-compose.*.yml | sort | xargs -I{} echo -f {}) up
```
