# La Familly Bass — Minimal Drive

This monorepo contains a very small Google-Drive-like service optimized for a Raspberry Pi 3 (ARM64, 1GB RAM).

Principles:
- Memory budget is critical: backend is a static Go binary; frontend served as static files by nginx.
- Media files stored on local SD in `./data/media` (mounted by docker-compose).
- SQLite used for metadata only (file `./data/db.sqlite`).

Quick start (local development):

1. Create data dirs:

```powershell
mkdir data\media -Force
```

2. Start services (Dev):

```powershell
docker-compose up --build
```

Defaults:
- Backend: http://localhost:8080
- Frontend: http://localhost (nginx on 80)

Environment variables (docker-compose):
- JWT_SECRET: secret for JWT
- ADMIN_EMAIL, ADMIN_PASSWORD: initial admin account created on boot

CI/CD:
- GitHub Actions workflow builds ARM64 images and pushes to Docker Hub. Set `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` in repository secrets.

Notes & optimizations for Raspberry Pi 3:
- Use `scratch` final image for the backend to minimize size.
- Keep logging minimal; consider disabling middleware.Logger in production to save CPU.
- Limit upload sizes in frontend and backend.
- Consider enabling swap file if SD is large enough, but avoid heavy swap thrashing.

