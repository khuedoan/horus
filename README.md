# Free Infrastructure

Always free (as in beer) multicloud infrastructure.
Can be used for personal use, student projects, or even start up PoC.

- An S3 bucket for Terraform backend
- Multicloud Kubernetes (k3s) cluster
  - Wireguard VPN
- Databases

## Get started

The default set up will creates all resources on all major cloud providers.

### Prerequisites

Create initial configurations, it will prompt you for the information.

```sh
make prerequisites
```

Continue reading to see what to fill in

#### Oracle Cloud

- Create an Oracle Cloud account
- Generate an API signing key:
  - Profile menu (User menu icon) -> User Settings -> API Keys -> Add API Key
  - Select Generate API Key Pair, download the private key and click Add
  - Check the Configuration File Preview for the values

Put those values in `./terraform.auto.tfvars`:

#### Google Cloud

TODO

#### Azure

TODO

#### AWS

TODO

### Provision

Create Terraform backend:

```sh
make backend
```

Build the infrastructure:

```sh
make init apply
```

## Usage

### VPN

Get QR code for mobile:

```sh
make vpn-config
```
