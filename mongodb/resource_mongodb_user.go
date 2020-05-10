package mongodb

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
)

func userResourceServer() *schema.Resource {
	return &schema.Resource{
		Create: ResourceMongoDBUserCreate,
		Read:   ResourceMongoDBUserRead,
		Update: ResourceMongoDBUserUpdate,
		Delete: ResourceMongoDBUserDelete,
		Exists: ResourceMongoDBUserExists,
		Importer: &schema.ResourceImporter{
			State: ResourceMongoDBUserImport,
		},

		CustomizeDiff: customdiff.All(
            customdiff.ForceNewIfChange("db", func (old, new, meta interface{}) bool {
                return new.(string) != old.(string)
            }),
		),

		Schema: map[string]*schema.Schema{
			"db": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"role": roleRefSet(),
		},
	}
}

func ResourceMongoDBUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	err := client.CreateUser(userInfo(d))

	if err != nil {
		return err
	}

	d.SetId(d.Get("username").(string))

	return nil
}

func ResourceMongoDBUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	database := d.Get("db").(string)

	obj, err := client.GetUser(database, d.Id())
	if err != nil {
		return err
	}

	if obj == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("role", flattenRoleRefs(obj.Roles)); err != nil {
		return err
	}

	// check password
	password := CheckUserPassword(database, d.Get("username").(string),d.Get("password").(string))
	if password != d.Get("password") {
		d.Set("password",password)
	}

	return nil
}

func ResourceMongoDBUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	if err := client.UpdateUser(userInfo(d)); err != nil {
		return err
	}

	return nil
}

func ResourceMongoDBUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	if err := client.DeleteUser(userInfo(d)); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func ResourceMongoDBUserExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*Client)
	database := d.Get("db").(string)

	obj, err := client.GetUser(database, d.Id())
	switch {
		case obj == nil:
			return false, nil
		case err != nil:
			return false, err
	}

	return true, nil
}

func ResourceMongoDBUserImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	client := m.(*Client)

	db, err := client.CheckUser(d.Id())
	if err != nil {
		return nil, err
	}

	obj, err := client.GetUser(db, d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("username", obj.Username)
	d.Set("db", obj.Db)

	if err := d.Set("role", flattenRoleRefs(obj.Roles)); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, err
}

func userInfo(d *schema.ResourceData) User {
	return User{
		Db:                          d.Get("db").(string),
		Description:                 d.Get("description").(string),
		Username:                    d.Get("username").(string),
		Password:                    d.Get("password").(string),
		Roles:                       expandRoleRefs(d.Get("role").(*schema.Set)),
	}
}
