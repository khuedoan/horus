.POSIX:
.EXPORT_ALL_VARIABLES:

KUBECONFIG = ../infra/kubeconfig.yaml

default: install

install:
	helm install bootstrap .

upgrade-flux:
	flux install --export > crds/gotk-components.yaml
