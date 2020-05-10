package mongodb

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"go.mongodb.org/mongo-driver/bson"
)

type RoleRef struct {
	Role string `bson:"role"`
	Db   string `bson:"db"`
}

func roleRefSet() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"db": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func flattenRoleRefs(in []RoleRef) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, 0)

	for _, v := range in {
		r := make(map[string]interface{})
		r["role"] = v.Role
		r["db"] = v.Db
		out = append(out, r)
	}

	return out
}

func expandRoleRefs(in *schema.Set) []RoleRef {
	out := make([]RoleRef, 0, 0)

	for _, v := range in.List() {
		r1 := v.(map[string]interface{})
		r := RoleRef{
			Role: r1["role"].(string),
			Db:   r1["db"].(string),
		}
		out = append(out, r)
	}

	return out
}

func bsonRoleRefs(roles []RoleRef) bson.A {
	out := bson.A{}

	for _, v := range roles {
		out = append(out, bson.M{
			"role": v.Role,
			"db":   v.Db,
		})
	}

	return out
}
