package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"errors"
)

type CheckRole struct {
	Role	string    `bson:"role"`
	Db		string    `bson:"db"`
}

type Role struct {
	Role       string      `bson:"role"`
	Db         string      `bson:"db"`
	Roles      []RoleRef   `bson:"roles"`
	Privileges []Privilege `bson:"privileges"`
}

type RoleInfoResult struct {
	Ok    int    `bson:"ok"`
	Roles []Role `bson:"roles"`
}

func (client *Client) GetRole(dbName string, role string) (*Role, error) {
	db := client.client.Database(dbName)

	command := bson.D{{"rolesInfo", role}, {"showPrivileges", true}}

	var result RoleInfoResult

	if err := db.RunCommand(*client.context, command).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Roles) == 1 {
		return &result.Roles[0], nil
	}

	return nil, nil
}

func (client *Client) CreateRole(role Role) error {
	db := client.client.Database(role.Db)

	command := bson.D{
		{"createRole", role.Role},
		{"roles", bsonRoleRefs(role.Roles)},
		{"privileges", bsonPrivileges(role.Privileges)},
	}

	var result CreateResult

	if err := db.RunCommand(*client.context, command).Decode(&result); err != nil {
		return err
	}

	return nil
}

func (client *Client) UpdateRole(role Role) error {
	db := client.client.Database(role.Db)

	command := bson.D{
		{"updateRole", role.Role},
		{"roles", bsonRoleRefs(role.Roles)},
		{"privileges", bsonPrivileges(role.Privileges)},
	}

	var result bson.M
	if err := db.RunCommand(*client.context, command).Decode(&result); err != nil {
		return err
	}

	return nil
}

func (client *Client) DeleteRole(role Role) error {
	db := client.client.Database(role.Db)

	command := bson.D{
		{"dropRole", role.Role},
	}

	var result bson.M
	if err := db.RunCommand(*client.context, command).Decode(&result); err != nil {
		return err
	}

	return nil
}

func (client *Client) CheckRole(name string) (string, error) {
	col := client.client.Database("admin").Collection("system.roles")

	var result CheckRole

	if err := col.FindOne(*client.context, bson.M{"role": name}).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Db) == 0 {
		log.Printf("[ERR] Could not find role [%s] in admin.system.roles", name)
		return "", errors.New("Could not find role in admin.system.roles")
	} else {
		log.Printf("[INFO] Found 1 result for role [%s]", name)
		return result.Db, nil
	}
}
