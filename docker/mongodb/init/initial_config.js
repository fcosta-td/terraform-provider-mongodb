db.getSiblingDB("db-test-1").createCollection("dummy")
db.getSiblingDB("db-test-2").createCollection("dummy")

// Example 2 - Import
db.getSiblingDB("db-test-2").createRole({ role: "role2", privileges: [], roles: [{ role: "read", db: "db-test-2" }] })
db.getSiblingDB("db-test-2").createUser({ user: "user2", pwd: "123", roles: [{ role: "role2", db: "db-test-2" }] })
