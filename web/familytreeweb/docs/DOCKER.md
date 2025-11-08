# Family Tree Web - Docker

This directory contains the Docker configuration for building and running the Family Tree Web UI.

## Features

- **Google Distroless Base**: Maximum security with minimal attack surface
- **Multi-stage Build**: Optimized image size
- **Non-root User**: Runs as UID 65532 for enhanced security
- **Static File Server**: Lightweight Go-based HTTP server
- **Security Headers**: X-Frame-Options, X-Content-Type-Options, etc.
- **Health Check**: Built-in health endpoint
- **Multi-arch Support**: AMD64 and ARM64

## Building

### Build for local testing

```bash
docker build -t familytreeweb:local .
```

### Build multi-arch

```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t familytreeweb:latest \
  --push \
  .
```

## Running

### Basic run

```bash
docker run -p 8080:8080 familytreeweb:local
```

### With environment variables

```bash
docker run -p 8080:8080 \
  -e VITE_KEYCLOAK_URL=http://localhost:15101 \
  -e VITE_KEYCLOAK_REALM=familytree \
  -e VITE_KEYCLOAK_CLIENT_ID=familytree-web \
  -e VITE_API_BASE_URL=http://localhost:15000 \
  familytreeweb:local
```

### With custom port

```bash
docker run -p 3000:8080 \
  -e PORT=8080 \
  familytreeweb:local
```

## Docker Compose

```yaml
version: '3.8'

services:
  familytreeweb:
    image: ghcr.io/rogerwesterbo/familytree-web:latest
    ports:
      - '8080:8080'
    environment:
      - VITE_KEYCLOAK_URL=http://keycloak:8080
      - VITE_KEYCLOAK_REALM=familytree
      - VITE_KEYCLOAK_CLIENT_ID=familytree-web
      - VITE_API_BASE_URL=http://familytree-api:8080
      - VITE_REDIRECT_URI=http://localhost:8080/callback
      - VITE_POST_LOGOUT_REDIRECT_URI=http://localhost:8080
    healthcheck:
      test: ['CMD', 'wget', '--quiet', '--tries=1', '--spider', 'http://localhost:8080/health']
      interval: 30s
      timeout: 3s
      retries: 3
    restart: unless-stopped
```

## Security

### Image Security Features

- ✅ **Distroless Base**: No shell, no package manager
- ✅ **Non-root User**: UID 65532
- ✅ **Read-only Root FS**: Enforced in Kubernetes
- ✅ **No Privileged Escalation**: Capabilities dropped
- ✅ **Security Scanning**: Trivy scans on every build
- ✅ **SBOM**: Software Bill of Materials included
- ✅ **Signed Images**: Cosign signatures

### Vulnerability Scanning

```bash
# Scan the image
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  aquasec/trivy image familytreeweb:local
```

## Development

### Local development (without Docker)

```bash
npm install
npm run dev
```

### Build for production

```bash
npm run build
```

### Preview production build

```bash
npm run preview
```

## Environment Variables

Build-time variables (set during `npm run build`):

- `VITE_KEYCLOAK_URL`: Keycloak server URL (default: `http://localhost:15101`)
- `VITE_KEYCLOAK_REALM`: Keycloak realm (default: `familytree`)
- `VITE_KEYCLOAK_CLIENT_ID`: Keycloak client ID (default: `familytree-web`)
- `VITE_REDIRECT_URI`: OAuth callback URI (default: `http://localhost:15200/callback`)
- `VITE_POST_LOGOUT_REDIRECT_URI`: Post-logout URI (default: `http://localhost:15200`)
- `VITE_API_BASE_URL`: Family Tree API URL (default: `http://localhost:15000`)

## Troubleshooting

### Container won't start

Check logs:

```bash
docker logs <container-id>
```

### Health check failing

Test health endpoint:

```bash
curl http://localhost:8080/health
```

### Permission denied errors

Ensure the image runs as non-root:

```bash
docker run --user 65532:65532 familytreeweb:local
```

## License

See the [LICENSE](../../LICENSE) file for details.
