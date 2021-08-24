# Free Infrastructure

Always free (as in beer) multicloud infrastructure.
Can be used for personal use, student projects, or even start up PoC.

We will mostly use Oracle Cloud because they have the most generous free tier.

- An S3 bucket for Terraform backend
- A Kubernetes cluster (4 nodes with 1 ARM core and 6 GB of memory)
- Databases:
  - Firestore
  - Azure Cosmos DB
- Bandwidth:
  - Unlimited inbound transfer
  - Inbound transfer (per month):
    - Oracle Cloud: 10 TB (yes, **Terabytes**)
    - Google Cloud: TODO
    - AWS: TODO
    - Azure: 5 GB

## Get started

The default set up will creates all resources.

### Prerequisites

- Fork this repository because you will need to customize it for your needs.
- A credit/debit card to register for the accounts.
- Basic knowledge on Terraform and Ansible (optional, but will help a lot)

Create initial configurations, it will prompt you for the information.

```sh
make prerequisites
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
