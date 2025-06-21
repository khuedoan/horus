# Apps

TODO automate this convention:

- Namespace is basically tenant
- 1 cluster per env (multi region workers, single region masters), so env is basically equivalant with cluster

`apps/$NAMESPACE/$APP/$CLUSTER.yaml`

```sh
export NAMESPACE=khuedoan
export APP=blog
export CLUSTER=local
helm template --namespace $NAMESPACE $APP oci://ghcr.io/bjw-s-labs/helm/app-template:4.1.1 --values $NAMESPACE/$APP/$CLUSTER.yaml > $CLUSTER.yaml
oras push docker.io/khuedoan/argocd-oci-demo-blog:$CLUSTER $CLUSTER.yaml
```
