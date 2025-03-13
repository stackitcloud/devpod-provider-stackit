# DevPod Provider STACKIT

<!-- markdownlint-disable-next-line MD013 MD034 -->
[![Go Report Card](https://goreportcard.com/badge/github.com/stackitcloud/devpod-provider-stackit)](https://goreportcard.com/report/github.com/stackitcloud/devpod-provider-stackit)

Run [DevPod](https://devpod.sh/) on [STACKIT](https://www.stackit.de).

## Usage

To use this provider in your DevPod setup, you will need to do the following steps:

1. See the [DevPod documentation](https://devpod.sh/docs/managing-providers/add-provider)
   for how to add a provider
2. Use the reference `stackitcloud/devpod-provider-stackit` to download the latest
   release from GitHub
3. Configure the provider by specifying a few options:
- Region
- Project ID
- Availability zone
- Disk size
- Private key path
- Service Account key path

### STACKIT Service Account

To authenticate with STACKIT you need to create a Service Account, generate a keypair locally and upload the public key
and download the Service Account Key JSON.
See the docs to see how: https://docs.stackit.cloud/stackit/en/create-a-service-account-key-175112456.html