.POSIX:

default: install

install:
	helm install bootstrap .

upgrade-flux:
	flux install --export > crds/gotk-components.yaml
