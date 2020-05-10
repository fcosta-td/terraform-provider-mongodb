# terraform-provider-mongodb

## Installation
Download source and build with `go build -o terraform-provider-mongodb`.
Move resulting binary to `~/.terraform.d/plugins` and then `chmod 0755 terraform-provider-mongodb`

## Example
See `examples` for usage examples.

## Features
- Roles
  - Role
  - Privileges
  - Terraform Import
- Users
  - User
  - Roles
  - Detects password changes
  - Terraform Import
