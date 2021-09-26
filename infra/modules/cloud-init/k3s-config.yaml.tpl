%{ if role == "server" }
# TODO multi server with embedded etcd in the same pool
cluster-init: true
disable-cloud-controller: true
disable:
- local-storage
- servicelb
- traefik
%{ endif }
token: ${token}
