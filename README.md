# speedtest

This is the cause:

![](Screenshot.png)

## Deployment

### With Docker

```bash
docker run -d -p 80:80 -e PORT=80 --name speedtest dihechen/speedtest:latest
```

### With Helm

1. Create a cert-manager Issuer.
2. Update `values.yaml` with your hosts or other options.

```bash
cd charts
helm install your-release speedtest
```
