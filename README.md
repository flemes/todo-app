# Todo App (Go REST API + React UI)

A small learning project: a Go REST service with in-memory storage and a separate React (Vite) UI.

```
todo-app/
├── backend/    Go 1.22 REST API (standard library only, no dependencies)
└── frontend/   React 18 + Vite single-page UI
```

## Run it

Requires Go 1.22+ and Node 18+. Use two terminals.

```powershell
# Terminal 1 - API on http://localhost:8080
cd backend
go run .

# Terminal 2 - UI on http://localhost:5173
cd frontend
npm install
npm run dev
```

Run the backend tests with `cd backend; go test ./...`.

## Configuration

| Where     | Variable         | Default                 | Purpose                              |
|-----------|------------------|-------------------------|--------------------------------------|
| backend   | `ADDR`           | `:8080`                 | Listen address                       |
| backend   | `ALLOWED_ORIGIN` | `http://localhost:5173` | Origin allowed by CORS (the UI)      |
| frontend  | `VITE_API_URL`   | `http://localhost:8080` | API base URL (see `.env.example`)    |

## API

JSON in/out. Errors look like `{"error": "message"}`.

| Method | Path              | Body                          | Success                  |
|--------|-------------------|-------------------------------|--------------------------|
| GET    | `/api/todos`      | -                             | `200` list of todos      |
| POST   | `/api/todos`      | `{"title": "..."}`            | `201` created todo       |
| GET    | `/api/todos/{id}` | -                             | `200` todo               |
| PUT    | `/api/todos/{id}` | `{"title": "...", "done": true}` | `200` updated todo    |
| DELETE | `/api/todos/{id}` | -                             | `204`                    |

Failures: `400` invalid id / JSON / empty title, `404` unknown id.

Todo shape: `{"id": 1, "title": "...", "done": false, "createdAt": "2026-01-01T00:00:00Z"}`

```powershell
curl.exe -X POST http://localhost:8080/api/todos -H "Content-Type: application/json" -d '{"title":"learn Go"}'
```

## Design notes

- **Storage** is an in-memory map guarded by a `sync.RWMutex` (`backend/store.go`). Data is lost on restart.
  To persist, keep the `Store` method set and swap the implementation (SQLite/Postgres) behind an interface.
- **Routing** uses the Go 1.22 `ServeMux` method+wildcard patterns (`GET /api/todos/{id}`), so no router library is needed.
- **CORS** middleware lets the separately-served UI call the API from the browser.
- **UI** keeps state locally and updates it from API responses; all HTTP calls live in `frontend/src/api.js`.

## Deliberately left out (future work)

Intentionally skipped to keep this a learning project. Revisit before any real deployment:

- **Authentication & authorization**: no login, no per-user todos; anyone who can reach the API can read/modify everything.
  Next step: JWT or session cookies (OIDC provider), a `userID` on each todo, and ownership checks in handlers.
- **Transport security**: plain HTTP only. Terminate TLS (reverse proxy or `ListenAndServeTLS`).
- **CORS**: currently a single allowed origin from config; keep it strict, never `*` once credentials exist.
- **Rate limiting / abuse protection**: none. Add per-IP limits and request timeouts
  (`http.Server` `ReadTimeout`, `WriteTimeout`, `IdleTimeout` - `http.ListenAndServe` has none).
- **Input limits**: body size is capped at 1 MiB, but title length is not limited. Add a max length.
- **Persistence & concurrency**: no database, migrations, or transactions.
- **Observability**: only a startup log line. Add request logging, metrics, and a health endpoint.
- **CSRF**: not applicable today (no cookies); becomes relevant if cookie-based sessions are added.
- **Frontend hardening**: security headers/CSP when serving the built UI; dependency audits (`npm audit`, `govulncheck`).
- **Tooling**: no CI, linting, Dockerfiles, or frontend tests yet.
