# Terraform Provider for Netcup CCP API

This is a [Terraform](https://terraform.io) provider for the [Netcup](https://www.netcup.de/) CCP [API](https://www.netcup-wiki.de/wiki/CCP_API). It provides resources to manage Netcup DNS records.

## Getting started
```terraform
terraform {
  required_providers {
    netcup-ccp = {
      source = "rincedd/netcup-ccp"
    }
  }
}

provider "netcup-ccp" {
  customer_number  = "123456"     # Netcup customer number
  ccp_api_key      = "xxxyyyzzz"  # API key for Netcup CCP
  ccp_api_password = "secret"     # API key password
}

resource "netcup-ccp_dns_record" "sample_record" {
  domain_name = "example.de"
  name        = "@"
  type        = "A"
  value       = "1.2.3.4"
  priority    = "0"               # for MX records
}

```
