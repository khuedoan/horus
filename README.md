# Cloudlab

> [!IMPORTANT]
> This project is designed to manage my offsite setup, which is specific to my
> use cases, so it might not be directly useful to you. For a ready-to-use
> solution, please refer to my [homelab project](https://github.com/khuedoan/homelab).

## Project structure

```
.
├── flake.nix                             # Contains dependencies required by this project for both local and CI/CD
├── Makefile                              # Entry point for all manual actions
├── compose.yaml                          # Servers required for running locally
├── infra                                 # Infrastructure definition
│   ├── modules                           # Terraform modules
│   │   ├── network
│   │   ├── instance
│   │   ├── cluster
│   │   └── ...
│   ├── local                             # Terragrunt configuration for the local environment
│   │   └── ...
│   └── ${ENV}                            # Terragrunt configuration for the ${ENV} environment
│       ├── root.hcl                      # Root config used by other Terragrunt files
│       ├── secrets.yaml                  # Encrypted secrets
│       ├── tfstate                       # Bootstrap Terraform state
│       ├── ${CLOUD}
│       │   └── ${REGION}
│       │       └── ${MODULE}
│       │           └── terragrunt.hcl
│       ├── metal
│       │   └── vn-south-1
│       │       ├── bootstrap
│       │       │   └── terragrunt.hcl
│       │       └── cluster
│       │           └── terragrunt.hcl
│       └── ...
├── platform                              # Highly privileged platform components
│   └── ${ENV}
│       ├── grafana.yaml
│       ├── temporal.yaml
│       ├── wireguard.yaml
│       └── ...
├── apps                                  # User applications, standardized with strict controls
│   ├── ${NAMESPACE}
│   │   └── ${APP}
│   │       └── ${ENV}.yaml
│   └── khuedoan
│       └── blog
│           ├── local.yaml
│           └── production.yaml
├── controller                            # Automation controller for the entire project - think GitHub Actions, but better
│   ├── activities                        # Temporal activities (git clone, terragrunt apply, etc.)
│   │   ├── git.go
│   │   ├── terragrunt.go
│   │   └── ...
│   ├── workflows                         # Temporal workflows, define a sequence of activities
│   │   ├── infra.go
│   │   ├── app.go
│   │   └── ...
│   ├── worker                            # Worker process that executes the workflows
│   └── Dockerfile                        # Builds the image for the controller, can run locally or on a cluster
└── test                                  # High level tests
```

## Features

### Infrastructure

Oracle Cloud was chosen due to their very generous free tier.

| Provider        | Service                     | Usage             | Pricing                        |
| :--             | :--                         | :--               | :--                            |
| Terraform Cloud | Workspace                   | 1                 | Free                           |
| Oracle Cloud    | Virtual Cloud Network       | 1                 | Free                           |
| Oracle Cloud    | `VM.Standard.A1.Flex` (ARM) | 4 cores, 24GB mem | Free (yes, you read it right!) |
| Oracle Cloud    | Block Storage               | 200GB             | Free                           |

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

Install the following packages:

- [Nix](https://nixos.org/download.html)

That's it! Run the following command to open the Nix shell:

```sh
nix develop
```

### Provision

Build the infrastructure:

```sh
make
```

## Acknowledgments and References

- [Oracle Terraform Modules](https://github.com/oracle-terraform-modules)
- [Official k3s systemd service file](https://github.com/k3s-io/k3s/blob/master/k3s.service)
- [Sample Prometheus configuration for Istio](https://github.com/istio/istio/blob/master/samples/addons/extras/prometheus-operator.yaml)
