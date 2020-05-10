provider "mongodb" {
  mongo_uri = "mongodb://localhost:27017/admin"
}

resource "mongodb_user" "user1" {
  db          = "db-test-1"
  description = "" # optional
  username    = "user1"
  password    = "123"

  role {
    role = mongodb_role.role1.name
    db   = mongodb_role.role1.db
  }
}

resource "mongodb_role" "role1" {
  name = "role1"
  db   = "admin"

  role {
    role = "read" // This is a built-in role in MongoDB
    db   = "db-test-1"
  }
  role {
    role = "readWrite" // This is a built-in role in MongoDB
    db   = "db-test-2"
  }

  privilege {
    cluster = true
    actions = ["listDatabases"]
  }

  privilege {
    db         = "admin"
    collection = "*"
    actions    = ["find"]
  }
}
