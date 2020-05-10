---
layout: "mongodb"
page_title: "Provider: MongoDB"
sidebar_current: "docs-mongodb-index"
description: |-
  A provider for MongoDB Server.
---

# MongoDB Provider

The MongoDB provider gives the ability to deploy and configure resources in a MongoDB server.

Use the navigation to the left to read about the available resources.

## Usage

```hcl
provider "mongodb" {
  mongo_uri = "mongodb://localhost:27017/admin"
}
```

```hcl
provider "postgresql" {
  mongo_uri = "mongodb://localhost:27017/admin"
}

resource "mongodb_user" "user1" {
  db          = "myDB"
  username    = "user1"
  password    = "supersecret"

  role {
    role = mongodb_role.role1.name
    db   = "admin"
  }
}

resource "mongodb_role" "role1" {
  name = "role1"
  db   = "admin"

  role {
    role = "read" // This is a built-in role in MongoDB
    db   = "myDB"
  }
}

```

## Argument Reference

The following arguments are supported:

* `mongo_uri` - (Required) MongoDB URI
