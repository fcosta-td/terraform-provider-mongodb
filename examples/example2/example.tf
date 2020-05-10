provider "mongodb" {
  mongo_uri = "mongodb://localhost:27017/admin"
}

resource "mongodb_user" "user2" {
  db          = "db-test-2"
  description = ""
  username    = "user2"
  password    = "123"

  role {
    role = "role2"
    db   = "db-test-2"
  }

  role {
    role = mongodb_role.role2.name
    db   = mongodb_role.role2.db
  }
}

resource "mongodb_role" "role2" {
  name = "role2"
  db   = "db-test-2"

  role {
    role = "read"
    db   = "db-test-2"
  }
}
