.POSIX:

apply:
	ansible-playbook \
		--inventory inventory.yml \
		--ask-vault-pass \
		main.yml

edit-vault:
	ansible-vault edit ./roles/global-secrets/vars/main.yml
