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

Create initial configurations, it will prompt you for the information.

```sh
make prereqs
```

Continue reading to see what to fill in.

<details> <summary>Oracle Cloud</summary>

- Create an Oracle Cloud account
- Generate an API signing key:
  - Profile menu (User menu icon) -> User Settings -> API Keys -> Add API Key
  - Select Generate API Key Pair, download the private key and click Add
  - Check the Configuration File Preview for the values

</details>

<details> <summary>Google Cloud</summary>

No resource on GCP yet

</details>

<details> <summary>AWS</summary>

No resource on AWS yet

</details>

<details> <summary>Azure</summary>

No resource on Azure yet

</details>

Remember to backup the credential files (you can put them in a password manager)

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
