### Log pipeline management for Kubernetes ###
Logging operator provides Kubernetes native log management for developers and devops teams. Main benefits are:

* Configure logging using Kubernetes constructs. No need to learn log configurations.
* Flexibility and Reuse through Kubernetes Custom Resource Definitions.
* Handles logging service deployment and scaling.
* Support for popular datastores like ElasticSearch and S3

#### Concepts ####
1. ***Output***: An output defines a datastore where logs are to be stored. Currently, the operator supports ElasticSearch and S3 as log stores.
Outputs are defined at cluster scope -- all logs from containers in the cluster get routed to each output.

#### Architecture ####
Logging operator uses fluent-bit and fluentd for collection and processing of logs respectively. The fluent-bit component is deployed as daemonset and is present on each node. Its main function is log formatting and filtering. fluentd is used as aggregator and buffer. It ships logs to chosen datastore. The fluentd layer can scale per log traffic.

![Architecture](docs/images/fluentd-arch.jpeg)


#### Install ####
Simplest way to install is with bundled deploy script
```
./hack/deploy.sh
```
If you are curious, deploy.sh creates prerequesite namespaces and applies yaml manifests under deploy/ directory.
#### Example Usage With Object Store ####
This example shows how to forward logs to an object storage. We will be using Minio, but the same example works with S3 as well.
0. Deploy Minio on Kubernetes
```bash
kubectl apply -f docs/getting-started/
```
This creates a secret containing api key and secret key referenced by minio process, a deployment for minio and a service to access minio api
1. Deploy nginx pod from Kubernetes example repo:
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/website/master/content/en/examples/application/deployment.yaml
```
2. Create a bucket named "test" in minio. Minio deployment creates a service named minio-service and listens on port 9000. To log into minio env, forward the local port 9000 to minio-service 9000 inside the cluster
```bash
kubectl port-forward svc/minio-service 9000
```
Then login to localhost:9000 from your browser.
The access key and secret key values can be extracted from secret "minio". Note that they are in base64 format. You need to decode before pasting values in the browser window.

3. We will now use this endpoint and credentials in its secret to configure a store for logs.
```yaml
apiVersion: logging.pf9.io/v1alpha1
kind: Output
metadata:
  name: objstore
spec:
  type: s3
  params:
    - name: aws_key_id
      valueFrom:
        secretKeyRef:
          name: minio
          key: MINIO_ACCESS_KEY
    - name: port
      value: 9000
    - name: aws_sec_key
      valueFrom:
        secretKeyRef:
          name: minio
          key: MINIO_SECRET_KEY
    - name: s3_bucket
      value: test
    - name: s3_region
      value: us-east-1 # default region used by minio
```

