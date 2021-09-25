#cloud-config

write_files:
  - path: /etc/rancher/k3s/config.yaml
    content: ${k3s_config}
    encoding: b64
    owner: root:root
    permission: '0600'
  - path: /etc/systemd/system/k3s.service
    content: ${k3s_service}
    encoding: b64
    owner: root:root
    permission: '0644'

runcmd:
  - curl https://github.com/k3s-io/k3s/releases/download/v1.21.4+k3s1/k3s -o /usr/local/bin/k3s
  - systemctl enable --now k3s
