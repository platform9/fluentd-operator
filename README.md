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
This example shows how to forward logs to an object storage. We will be using S3.
1. Deploy random logger as kubernetes deployment
```bash
kubectl apply -f docs/getting-started/user-guides/random-logger.yaml
```
2. Create a bucket named "test" in S3.
3. Create a kubernetes secret and specify your access and secret keys to gain access to this S3 bucket
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: s3
type: Opaque
stringData:
  access_key: <aws access key>
  secret_key: <aws secret key>
```
4. Create output CR indicating S3 as destination. Note the AWS secrets referenced by specifying the secret name and key within it.
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
        name: s3
        namespace: default
        key: access_key
    - name: aws_sec_key
      valueFrom:
        name: s3
        namespace: default
        key: secret_key
    - name: s3_region
      value: <s3 region name>
    - name: s3_bucket
      value: <s3 bucket name>
```

