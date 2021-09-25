%{ if role == "server" }
# TODO multi master with embedded etcd in the same pool
cluster-init: true
disable:
- local-storage
- traefik
%{ endif }
token: ${token}
