# MurmurationsAllocator
## Create secrets
1. Create secret for MurmurationsAllocator
```
kubectl create secret generic allocator-app-secret \
  --from-literal="MONGO_HOST=mongodb+srv://" \
  --from-literal="MONGO_USERNAME=mongo-admin" \
  --from-literal="MONGO_PASSWORD=mongo-password"
```