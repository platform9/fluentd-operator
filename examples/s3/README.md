## Fluentd operator with S3

This example shows how to forward logs to an s3 object storage.

1. Create S3 bucket if not already created and note the access-key and secret-key. 
2. Deploy random logger as kubernetes deployment
```bash
kubectl apply -f ../../docs/getting-started/user-guides/random-logger.yaml
```
3. Create a kubernetes secret and specify your access and secret keys to gain access to this S3 bucket
```bash
kubectl apply -f templates/secret-s3.yaml
```
4. Create output CR indicating S3 bucket as destination. Note the AWS secrets referenced by specifying the secret name and key within it
```bash
kubectl apply -f templates/object-s3.yaml
```