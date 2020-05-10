Example 2 - Import
===========
This example tests terraform resource import functionality for previously created user2 and role2 (already available on docker-compose).


## Import role

```hcl
terraform import mongodb_role.role2 role2
```

## Import user

```hcl
terraform import mongodb_user.user2 user2
```
