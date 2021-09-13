# Free Infrastructure

```diff
! ⚠️ WORK IN PROGRESS
```

Always free (as in beer) cloud infrastructure.

This repo is meant to be forked to customize for your needs, it can be used for self-hosting, hobby projects, student assignments...

## Features

### Infrastructure

- 1 VM for a VPN server
- 1 VM for a mail server
- 1 Kubernetes cluser

## Applications

- Wireguard VPN server
- Mailcow mail server
- CI/CD with Tekton and ArgoCD
- TBD

### Free tier limits

The following list only includes the services that we use in this repository.

| Cloud | Name | Purpose | Using | Limit | Notes |
| ----- | ---- | ------- | ----- | ----- | ----- |
| AWS | S3 | Terraform backend | 1 | 1 | Normal Terraform usage should not exceed the capacity or bandwidth limit |
| Oracle Cloud | VM (x86) | Mail and VPN server | 2 | 2 | None |
| Oracle Cloud | VM (ARM) | Kubernetes cluster | 4 | 4 | None |

## Get started

### Prerequisites

- Fork this repository because you will need to customize it for your needs.
- A credit/debit card to register for the accounts.
- Basic knowledge on Terraform and Ansible (optional, but will help a lot)

Configuration files:

<details>

<summary>Terraform Cloud</summary>

- Create a Terraform Cloud account at <https://app.terraform.io>
- Run `terraform login` and follow the instruction

</details>

<details>

<summary>Oracle Cloud</summary>

- Create an Oracle Cloud account at <https://cloud.oracle.com>
- Generate an API signing key:
  - Profile menu (User menu icon) -> User Settings -> API Keys -> Add API Key
  - Select Generate API Key Pair, download the private key to `~/.oci/private.pem` and click Add
  - Copy the Configuration File Preview to `~/.oci/config` and change `key_file` to `~/.oci/private.pem`

</details>

Remember to backup the following credential files (you can put them in a password manager):

- `~/.terraform.d/credentials.tfrc.json`
- `~/.oci/config`
- `~/.oci/private.pem`

### Provision

Build the infrastructure:

```sh
make infra
```

## Usage

### VPN

Get QR code for mobile:

```sh
./scripts/vpn-config
```
