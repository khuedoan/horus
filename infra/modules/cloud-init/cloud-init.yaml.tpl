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
  - netfilter-persistent flush
  - systemctl disable --now netfilter-persistent
  - curl -L "https://github.com/k3s-io/k3s/releases/download/v1.21.4+k3s1/k3s$(case $(uname -m) in arm64|aarch64) echo '-arm64' ;; amd64|x86_64) echo '' ;; esac)" -o /usr/local/bin/k3s
  - chmod +x /usr/local/bin/k3s
  - ln -s /usr/local/bin/k3s /usr/local/bin/kubectl
  - systemctl enable --now k3s
