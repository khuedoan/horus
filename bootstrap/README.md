# Bootstrap

## Generate Flux manifests

```sh
flux install --export > flux-system/gotk-components.yaml
```

## Install Flux and bootstrap components

```sh
kustomize build flux-system | kubectl apply -f -
```
