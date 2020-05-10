package mongodb

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
)

func roleResourceServer() *schema.Resource {
	return &schema.Resource{
		Create: ResourceMongoDBRoleCreate,
		Read:   ResourceMongoDBRoleRead,
		Update: ResourceMongoDBRoleUpdate,
		Delete: ResourceMongoDBRoleDelete,
		Exists: ResourceMongoDBRoleExists,
		Importer: &schema.ResourceImporter{
			State: ResourceMongoDBRoleImport,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role":      roleRefSet(),
			"privilege": privilegeSet(),
		},
	}
}

func ResourceMongoDBRoleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	if err := client.CreateRole(roleInfo(d)); err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))

	return nil
}

func ResourceMongoDBRoleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	obj, err := client.GetRole(d.Get("db").(string), d.Id())
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

	if err := d.Set("privilege", flattenPrivileges(obj.Privileges)); err != nil {
		return err
	}

	return nil
}

func ResourceMongoDBRoleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	if err := client.UpdateRole(roleInfo(d)); err != nil {
		return err
	}

	return nil
}

func ResourceMongoDBRoleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	if err := client.DeleteRole(roleInfo(d)); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func ResourceMongoDBRoleExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*Client)
	database := d.Get("db").(string)

	obj, err := client.GetRole(database, d.Id())
	switch {
		case obj == nil:
			return false, nil
		case err != nil:
			return false, err
	}

	return true, nil
}

func ResourceMongoDBRoleImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	client := m.(*Client)

	db, err := client.CheckRole(d.Id())
	if err != nil {
		return nil, err
	}

	obj, err := client.GetRole(db, d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("name", obj.Role)
	d.Set("db", obj.Db)

	if err := d.Set("role", flattenRoleRefs(obj.Roles)); err != nil {
		return nil, err
	}

	if err := d.Set("privilege", flattenPrivileges(obj.Privileges)); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, err
}

func roleInfo(d *schema.ResourceData) Role {
	return Role{
		Db:         d.Get("db").(string),
		Role:       d.Get("name").(string),
		Roles:      expandRoleRefs(d.Get("role").(*schema.Set)),
		Privileges: expandPrivileges(d.Get("privilege").(*schema.Set)),
	}
}
