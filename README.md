# MurmurationsAllocator
## Run locally
1. Install NGINX Ingress Controller
   ```
   helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
    helm repo update
    helm upgrade --install ingress-nginx ingress-nginx \
    --repo https://kubernetes.github.io/ingress-nginx \
    --namespace ingress-nginx --create-namespace
   ```
2. Create secret for MurmurationsAllocator
    ```
    kubectl create secret generic allocator-app-secret \
      --from-literal="MONGO_HOST=mongodb+srv://" \
      --from-literal="MONGO_USERNAME=mongo-admin" \
      --from-literal="MONGO_PASSWORD=mongo-password"
    ```
3. Add the following to your host file sudo vim /etc/hosts
   ```
   127.0.0.1    allocator.murmurations.dev
   ```
