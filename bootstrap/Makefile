.POSIX:

default: apply

apply:
	kustomize build flux-system | kubectl apply -f -

upgrade:
	flux install --export > flux-system/gotk-components.yaml
