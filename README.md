# DevPod Provider for STACKIT

<!-- markdownlint-disable-next-line MD013 MD034 -->
[![Open in DevPod](https://devpod.sh/assets/open-in-devpod.svg)](https://devpod.sh/open#provider=stackitcloud/devpod-provider-stackit)

[![Go Report Card](https://goreportcard.com/badge/github.com/stackitcloud/devpod-provider-stackit)](https://goreportcard.com/report/github.com/stackitcloud/devpod-provider-stackit)

The DevPod Provider for STACKIT enables seamless integration between [DevPod](https://devpod.sh) and the [STACKIT Cloud Platform](https://www.stackit.de/en/), allowing you to deploy development environments on STACKIT's European cloud infrastructure. This provider offers a secure, GDPR-compliant solution for creating reproducible development environments in STACKIT's German data centers.

## Features

- Deploy development environments on STACKIT's cloud platform
- Automated VM provisioning and management
- GDPR-compliant European cloud infrastructure
- Full DevContainer compatibility

## Requirements

- [DevPod CLI](https://devpod.sh/docs/getting-started/installation) or [DevPod Desktop](https://devpod.sh/docs/getting-started/installation)
- STACKIT account with an active project
- Service Account with appropriate permissions
- Service Account Key and Private Key (for authentication)

## Installation

You can add the STACKIT provider to DevPod using the CLI:

```bash
devpod provider add stackitcloud/devpod-provider-stackit
```

Or through the DevPod Desktop application by clicking on "Providers" and then "Add Provider".

## Configuration

When adding the provider, you need to configure the following options:

| Option | Required | Description | Default  |
|--------|----------|-------------|----------|
| REGION | true | The STACKIT region to create the VM in (e.g., `eu01`) | `eu01`   |
| AVAILABILITY_ZONE | true | The availability zone to use (e.g., `eu01-1`, `eu01-2`) |    `eu01-1`     |
| PROJECT_ID | true | The STACKIT project ID to use |          |
| DISK_SIZE | false | The disk size to use in GB | `64`     |
| SERVICE_ACCOUNT_KEY_PATH | true | Path to your STACKIT Service Account Key JSON file |          |
| PRIVATE_KEY_PATH | true | Path to your private key |          |
| FLAVOR | false | The VM instance type to use | `g1.1`   |

## Authentication

To authenticate with STACKIT, you need to:

1. Create a Service Account in the STACKIT Portal
2. Generate a key pair locally
3. Upload the public key to STACKIT
4. Download the Service Account Key JSON file

### Creating a Service Account Key

1. Log in to the [STACKIT Portal](https://portal.stackit.cloud/)
2. Navigate to your project
3. Go to "Service Accounts" tab
4. Create a new Service Account or select an existing one
5. Go to "Service Account Keys" and create a new key
6. Follow the instructions to upload your public key and download the Service Account Key JSON

For detailed instructions, see the [STACKIT documentation on creating service account keys](https://docs.stackit.cloud/stackit/en/create-a-service-account-key-175112456.html).

## Usage

After setting up the provider, you can create a new DevPod workspace:

```bash
# Create a new workspace from a Git repository
devpod up git@github.com:user/repository.git --provider stackit

# Or from a local directory
devpod up ./my-project --provider stackit
```


## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the [Apache License 2.0](LICENSE).