# Kubernetes Deployment for GoDE

This directory contains Kubernetes manifests for deploying GoDE to a local minikube cluster for development and testing.

## Prerequisites

- [minikube](https://minikube.sigs.k8s.io/docs/start/) v1.30+
- [kubectl](https://kubernetes.io/docs/tasks/tools/) v1.28+
- [Docker](https://docs.docker.com/get-docker/) (for building images)

## Quick Start

### 1. Start Minikube

```bash
minikube start --cpus=4 --memory=4096
```

### 2. Build Docker Image

Build the Docker image in minikube's Docker environment:

```bash
# Point your shell to minikube's docker-daemon
eval $(minikube docker-env)

# Build the image
docker build -t gode-server:latest .
```

### 3. Deploy to Kubernetes

Apply all manifests in order:

```bash
# Create ConfigMap and Secret
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml

# Deploy database and cache layers
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/redis.yaml

# Wait for databases to be ready
kubectl wait --for=condition=ready pod -l app=postgres --timeout=120s
kubectl wait --for=condition=ready pod -l app=redis --timeout=60s

# Deploy the application
kubectl apply -f k8s/deserver.yaml

# Wait for application to be ready
kubectl wait --for=condition=ready pod -l app=deserver --timeout=120s
```

### 4. Access the Application

Get the NodePort for HTTP access:

```bash
kubectl get svc deserver-http
```

Access the service:

```bash
# Get minikube IP
minikube ip

# Access the HTTP gateway (replace NODE_PORT with the port from above)
curl http://$(minikube ip):NODE_PORT/health
```

Or use minikube service:

```bash
minikube service deserver-http
```

## Architecture

The deployment consists of:

- **PostgreSQL**: Persistent storage for executions, users, and Pareto sets
- **Redis**: Cache layer for fast execution lookups and pub/sub for real-time updates
- **deserver**: GoDE application server (2 replicas for HA)

### Resource Allocation

| Component | CPU Request | CPU Limit | Memory Request | Memory Limit |
|-----------|-------------|-----------|----------------|--------------|
| PostgreSQL | 100m | 500m | 128Mi | 512Mi |
| Redis | 50m | 250m | 64Mi | 256Mi |
| deserver | 200m | 1000m | 256Mi | 1Gi |

## Configuration

### Environment Variables

Configuration is managed through ConfigMap (`k8s/configmap.yaml`) and Secret (`k8s/secret.yaml`).

**ConfigMap** (non-sensitive):
- Server ports (gRPC: 3030, HTTP: 8081)
- Database connection details
- Redis connection details
- TTL settings for executions and progress
- Circuit breaker configuration

**Secret** (sensitive - change in production):
- `DB_PASSWORD`: PostgreSQL password
- `REDIS_PASSWORD`: Redis password (empty for dev)
- `JWT_SECRET`: JWT signing secret

### Customization

To modify configuration:

1. Edit `k8s/configmap.yaml` or `k8s/secret.yaml`
2. Apply changes: `kubectl apply -f k8s/configmap.yaml`
3. Restart pods to pick up changes: `kubectl rollout restart deployment/deserver`

## Health Checks

The application exposes health endpoints:

- `/health`: Liveness probe (application is running)
- `/ready`: Readiness probe (application is ready to serve traffic)

Check pod health:

```bash
kubectl get pods -l app=deserver
kubectl describe pod <pod-name>
```

## Logs

View application logs:

```bash
# All deserver pods
kubectl logs -l app=deserver -f

# Specific pod
kubectl logs <pod-name> -f

# Previous instance (if pod crashed)
kubectl logs <pod-name> -p
```

## Scaling

Scale the application:

```bash
# Scale to 3 replicas
kubectl scale deployment deserver --replicas=3

# Verify
kubectl get deployment deserver
```

## Troubleshooting

### Pods not starting

Check pod status and events:

```bash
kubectl get pods
kubectl describe pod <pod-name>
```

### Database connection issues

Verify PostgreSQL is ready:

```bash
kubectl exec -it <postgres-pod-name> -- psql -U gode -d gode -c "SELECT 1"
```

### Redis connection issues

Verify Redis is ready:

```bash
kubectl exec -it <redis-pod-name> -- redis-cli ping
```

### Port forwarding for debugging

Forward ports to localhost:

```bash
# gRPC port
kubectl port-forward svc/deserver-grpc 3030:3030

# HTTP port
kubectl port-forward svc/deserver-http 8081:8081

# PostgreSQL (for debugging)
kubectl port-forward svc/postgres-service 5432:5432

# Redis (for debugging)
kubectl port-forward svc/redis-service 6379:6379
```

## Cleanup

Remove all resources:

```bash
kubectl delete -f k8s/deserver.yaml
kubectl delete -f k8s/redis.yaml
kubectl delete -f k8s/postgres.yaml
kubectl delete -f k8s/secret.yaml
kubectl delete -f k8s/configmap.yaml

# Remove PVC (persistent data will be lost)
kubectl delete pvc postgres-pvc
```

Stop minikube:

```bash
minikube stop
# Or delete the cluster entirely
minikube delete
```

## Production Considerations

This configuration is optimized for development with minikube. For production:

1. **Security**:
   - Change all passwords in `k8s/secret.yaml`
   - Use a proper secret management solution (e.g., HashiCorp Vault, Sealed Secrets)
   - Enable TLS for gRPC and HTTPS for the gateway
   - Use network policies to restrict pod communication

2. **High Availability**:
   - Deploy PostgreSQL with replication (e.g., using CloudNativePG operator)
   - Deploy Redis with sentinel or cluster mode
   - Increase deserver replicas (3+ for production)
   - Use pod anti-affinity rules

3. **Storage**:
   - Use production-grade storage class (not hostPath)
   - Configure backup strategies for PostgreSQL
   - Consider StatefulSets for stateful components

4. **Monitoring**:
   - Add Prometheus metrics endpoints
   - Deploy Grafana dashboards
   - Configure alerting rules
   - Implement distributed tracing

5. **Resource Management**:
   - Set appropriate resource limits based on load testing
   - Configure horizontal pod autoscaling (HPA)
   - Use pod disruption budgets (PDB)

6. **Ingress**:
   - Replace NodePort with proper Ingress controller
   - Configure domain names and TLS certificates
   - Implement rate limiting and WAF rules
