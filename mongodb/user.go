package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"errors"
)

type CheckUser struct {
	Username	string    `bson:"user"`
	Db			string    `bson:"db"`
}

type User struct {
	Username		string    `bson:"user"`
	Password		string    `bson:"-"`
	Description		string    `bson:"customData.description"`
	Db				string    `bson:"db"`
	Roles			[]RoleRef `bson:"roles"`
}

type UsersInfoResult struct {
	Ok    int    `bson:"ok"`
	Users []User `bson:"users"`
}

func (client *Client) GetUser(databaseName string, name string) (*User, error) {
	db := client.client.Database(databaseName)

	command := bson.D{{"usersInfo", name}}

	var result UsersInfoResult

	if err := db.RunCommand(*client.context, command).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Users) == 1 {
		return &result.Users[0], nil
	}

	return nil, nil
}

func (client *Client) CreateUser(user User) error {
	db := client.client.Database(user.Db)

	command := bson.D{
		{"createUser", user.Username},
		{"pwd", user.Password},
		{"customData", bson.D{{"description", user.Description}, {"managedByTerraform", true}}},
		{"roles", bsonRoleRefs(user.Roles)},
	}

	var result CreateResult

	if err := db.RunCommand(*client.context, command).Decode(&result); err != nil {
		return err
	}

	return nil
}

func (client *Client) UpdateUser(user User) error {
	db := client.client.Database(user.Db)

	command := bson.D{
		{"updateUser", user.Username},
		{"customData", bson.D{{"description", user.Description}, {"managedByTerraform", true}}},
		{"roles", bsonRoleRefs(user.Roles)},
		{"pwd", user.Password},
	}

	var result CreateResult

	if err := db.RunCommand(*client.context, command).Decode(&result); err != nil {
		return err
	}

	return nil
}

func (client *Client) DeleteUser(user User) error {
	db := client.client.Database(user.Db)

	command := bson.D{{"dropUser", user.Username}}

	var result CreateResult

	if err := db.RunCommand(*client.context, command).Decode(&result); err != nil {
		return err
	}

	return nil
}

func (client *Client) CheckUser(name string) (string, error) {
	col := client.client.Database("admin").Collection("system.users")

	var result CheckUser

	if err := col.FindOne(*client.context, bson.M{"user": name}).Decode(&result); err != nil {
		return "", err
	}

	switch len(result.Db) {
		case 0:
			log.Printf("[ERR] Could not find user [%s] in admin.system.users", name)
			return "", errors.New("Could not find user in admin.system.users")

		case 1:
			log.Printf("[INFO] Found 1 result for user [%s]", name)
			return result.Db, nil

		default:
			log.Printf("[ERR] Found (%q) matches for user [%s] in admin.system.users",len(result.Db), name)
			return "", errors.New("Found more that one match for user in admin.system.users")
	}
}
