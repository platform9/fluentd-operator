## Fluentd operator with elastic search

This example shows how to forward logs to an elastic search object storage.

1. Check if elastic search is running and note the URL. If elastic search not installed, run below command to setup minimal elastic search deployment
```bash
kubectl apply -f templates/deploy-elasticsearch.yaml
```
2. Deploy random logger as kubernetes deployment
```bash
kubectl apply -f ../../docs/getting-started/user-guides/random-logger.yaml
```
3. Create output CR indicating elasticsearch as destination. Note the elasticsearch URL referenced in the params.
```yaml
apiVersion: logging.pf9.io/v1alpha1
kind: Output
metadata:
  name: objstore
spec:
  type: elasticsearch
  params:
    - name: url
      value: <elastic-search url>
    - name: user
      value: <user-name>
    - name: password
      value: <user-password>
    - name: index_name
      value: <index-name>
```
