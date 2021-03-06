# Terraform Guard

![Go](https://github.com/skillingbeck/tfguard/workflows/Go/badge.svg)

A command line tool to prevent unintented destruction of terraform resources.

Intended to be run between `terraform plan` and `terraform apply` when using CICD pipelines to deploy infrastructure changes. If the plan includes changes that involve a `destroy` action on any resource that isn't explicitly on an allow list, then the exit code will be non-zero.

With Terraform Guard in the CICD pipeline, anyone making infrastructure changes is forced to signal their intention for resources to be removed or replaced, thereby preventing disruption due to accidental destruction of resources. It also supports the Expand-Contract pattern for zero-downtime deployments:

1. Create replacement resource in terraform alongside original resource; deploy changes
2. Redirect traffic/associations from the original resource to the replacement
3. Verify everything is working as intended and original resource is no longer in use
4. Explicitly allow original resource to be deleted in `tfguard`; remove original resource from terraform; deploy changes

Supports terraform 0.12 and above.

## Example

```
terraform plan -out=out.tfplan -input=false
terraform show -json out.tfplan > ../tfplan.json
tfguard tfplan.json
terraform apply out.tfplan -auto-approve -input=false
```

## Configuration

### Allow certain resource(s) to be destroyed at a given address
A comma separated list of addresses. Partial paths are permitted, which will allow any resources that start with that path to be destroyed.

```
tfguard -allow-address-destroy=address1,address2
```

### Allow certain resource types to be destroyed
A comma separated list of types.

```
tfguard -allow-type-destroy=type1,type2
```
