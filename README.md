# Free Infrastructure

```diff
! ⚠️ WORK IN PROGRESS
```

Always free (as in beer) cloud infrastructure.

This repo is meant to be forked to customize for your needs, it can be used for self-hosting, hobby projects, student assignments...

## Features

### Infrastructure

| Provider        | Service                        | Using             | Limit             | Notes                 |
|-----------------|--------------------------------|-------------------|-------------------|-----------------------|
| Terraform Cloud | Workspace                      | 1                 | None              |                       |
| Oracle Cloud    | `VM.Standard.E2.1.Micro` (x86) | 2 instances       | 2 instances       |                       |
| Oracle Cloud    | `VM.Standard.A1.Flex` (ARM)    | 4 cores, 24GB mem | 4 cores, 24GB mem | 2x(2 cores, 12GB mem) |
| Oracle Cloud    | Block Storage                  | 200GB             | 200GB             | 50GB per VM           |
| Oracle Cloud    | Virtual Cloud Network          | 1                 | 2                 |                       |

### Applications

- Wireguard VPN server
- Mailcow mail server
- CI/CD with Tekton and ArgoCD
- TBD

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

If you see a warning like this, try to avoid those regions:

> ⚠️ Because of high demand for Arm Ampere A1 Compute capacity in the Foo and Bar regions, A1 instance availability in these regions is limited.
> If you plan to create A1 instances, we recommend choosing another region as your home region

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

## Acknowledgments and References

- [Oracle Terraform Modules](https://github.com/oracle-terraform-modules)
