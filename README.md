# Horus

> [!IMPORTANT]
> This project is designed to manage my offsite setup, which is specific to my
> use cases, so it might not be directly useful to you. For a ready-to-use
> solution, please refer to my [homelab project](https://github.com/khuedoan/homelab).

> The name is from [Horus the Child, or Harpocrates](https://en.wikipedia.org/wiki/Harpocrates)

## Features

### Infrastructure

Oracle Cloud was chosen due to their very generous free tier.

| Provider        | Service                     | Usage             | Pricing                        |
| :--             | :--                         | :--               | :--                            |
| Terraform Cloud | Workspace                   | 1                 | Free                           |
| Oracle Cloud    | Virtual Cloud Network       | 1                 | Free                           |
| Oracle Cloud    | `VM.Standard.A1.Flex` (ARM) | 4 cores, 24GB mem | Free (yes, you read it right!) |
| Oracle Cloud    | Block Storage               | 200GB             | Free                           |

### Applications

- [x] Hugo blog
- [x] Mail server (receive only, [outbound SMTP is blocked by default](https://docs.oracle.com/en-us/iaas/releasenotes/changes/f7e95770-9844-43db-916c-6ccbaf2cfe24/), you must upgrade to a paid account to open a service limits request to request an exemption)
- [ ] Web analytics
- [ ] Wireguard VPN server
- [ ] TBD

## Get started

### Prerequisites

- Fork this repository because you will need to customize it for your needs.
- A credit/debit card to register for the accounts.
- Basic knowledge on Terraform, Ansible and Kubernetes (optional, but will help a lot)

Configuration files:

<details>

<summary>Terraform Cloud</summary>

- Create a Terraform Cloud account at <https://app.terraform.io>

</details>

<details>

<summary>Oracle Cloud</summary>

- Create an Oracle Cloud account at <https://cloud.oracle.com>
- Generate an API signing key:
  - Profile menu (User menu icon) -> User Settings -> API Keys -> Add API Key
  - Select Generate API Key Pair, download the private key to `~/.oci/private.pem` and click Add
  - Copy the Configuration File Preview to `~/.oci/config` and change `key_file` to `~/.oci/private.pem`

If you see a warning like this, try to avoid those regions:

> ⚠️ Because of high demand for Arm Ampere A1 Compute capacity in the Foo and Bar regions, A1 instance availability in these regions is limited.
> If you plan to create A1 instances, we recommend choosing another region as your home region

</details>

Remember to backup the following credential files (you can put them in a password manager):

- `~/.terraform.d/credentials.tfrc.json`

Install the following packages:

- [Nix](https://nixos.org/download.html)

That's it! Run the following command to open the Nix shell:

```sh
nix-shell
```

### Provision

Build the infrastructure:

```sh
make
```

## Usage

### Connect to Wireguard

Desktop:

```sh
# kubectl exec to the wireguard pod
cat /config/peer_desktop/peer_desktop.conf
# Copy the content to horus.conf
# Then import the config file, for example on Linux using NetworkManager
nmcli connection import type wireguard file horus.conf
nmcli connection up horus
```

Mobile:

```sh
# kubectl exec to the wireguard pod
/app/show-peer phone
# Then scan the QR code
```

### Other

TBD

## Acknowledgments and References

- [Oracle Terraform Modules](https://github.com/oracle-terraform-modules)
- [Official k3s systemd service file](https://github.com/k3s-io/k3s/blob/master/k3s.service)
