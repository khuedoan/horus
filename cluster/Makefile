.POSIX:

IP_ADDRESS = $(shell cd ../infra && terraform output -raw instance_public_ip)

apply:
	ansible-playbook \
		--inventory "${IP_ADDRESS}", \
		main.yml
