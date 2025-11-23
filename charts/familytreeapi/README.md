# FamilyTree API Helm Chart

This Helm chart deploys the FamilyTree API application to Kubernetes.

## HTTPS/TLS Configuration

The API supports HTTPS using TLS certificates. To enable HTTPS in Kubernetes:

### 1. Create TLS Secret

First, create a Kubernetes secret with your TLS certificate and key:

```bash
kubectl create secret tls familytreeapi-tls \
  --cert=path/to/tls.crt \
  --key=path/to/tls.key \
  --namespace=your-namespace
```

### 2. Enable HTTPS in values.yaml

```yaml
httpApi:
  secure: true
  tlsSecretName: familytreeapi-tls
  corsAllowedOrigins: "https://example.com,https://app.example.com"
```

### 3. Deploy or Upgrade

```bash
helm upgrade --install familytreeapi ./charts/familytreeapi \
  --namespace your-namespace \
  --values custom-values.yaml
```

## Local Development

For local development, HTTPS is enabled by default using certificates from `hack/certs/`.

### Enable HTTPS (default)
```bash
# In .env
HTTP_API_SECURE=true
HTTP_API_TLS_CERT_FILE=hack/certs/server.crt
HTTP_API_TLS_KEY_FILE=hack/certs/server.key
```

### Disable HTTPS
```bash
# In .env
HTTP_API_SECURE=false
```

Or in VS Code launch.json:
```json
{
  "env": {
    "HTTP_API_SECURE": "false"
  }
}
```

## Configuration Options

| Parameter | Description | Default |
|-----------|-------------|---------|
| `httpApi.secure` | Enable HTTPS/TLS | `false` |
| `httpApi.tlsSecretName` | Name of Kubernetes secret containing TLS certs | `""` |
| `httpApi.corsAllowedOrigins` | Comma-separated list of allowed CORS origins | `""` |

## Notes

- The API uses RS512 algorithm for JWT token verification with Keycloak
- Self-signed certificates will show browser warnings in development
- In production, use proper CA-signed certificates
