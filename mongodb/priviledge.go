package mongodb

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"go.mongodb.org/mongo-driver/bson"
)

type Privilege struct {
	Resource Resource `bson:"resource"`
	Actions  []string `bson:"actions"`
}

type Resource struct {
	Cluster    bool    `bson:"cluster"`
	Db         *string `bson:"db,omitempty"`
	Collection *string `bson:"collection,omitempty"`
}

// API -> Terraform
func flattenPrivileges(in []Privilege) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, 0)

	for _, v := range in {
		p := make(map[string]interface{})
		if v.Resource.Cluster {
			p["cluster"] = true
			p["db"] = nil
			p["collection"] = nil
		} else {
			p["db"] = v.Resource.Db
			if *v.Resource.Collection == "" {
				p["collection"] = "*"
			} else {
				p["collection"] = v.Resource.Collection
			}
		}
		p["actions"] = v.Actions
		out = append(out, p)
	}

	return out
}

// Terraform -> API
func expandPrivileges(in *schema.Set) []Privilege {
	out := make([]Privilege, 0, 0)
	for _, v := range in.List() {
		p := v.(map[string]interface{})

		actions := make([]string, 0, 0)
		for _, a := range p["actions"].(*schema.Set).List() {
			actions = append(actions, a.(string))
		}

		cluster := p["cluster"].(bool)
		db := p["db"].(string)
		collection := p["collection"].(string)

		r := Resource{}
		if !cluster {
			r.Cluster = false
			r.Db = &db
			r.Collection = &collection
		} else {
			r.Cluster = true
			r.Db = nil
			r.Collection = nil
		}

		out = append(out, Privilege{
			Resource: r,
			Actions:  actions,
		})
	}
	return out
}

func bsonPrivileges(privileges []Privilege) bson.A {
	out := bson.A{}

	for _, v := range privileges {
		actions := bson.A{}
		for _, a := range v.Actions {
			actions = append(actions, a)
		}
		resource := bson.M{}
		if v.Resource.Cluster {
			resource["cluster"] = true
		} else {
			resource["db"] = *v.Resource.Db
			resource["collection"] = *v.Resource.Collection
		}
		out = append(out, bson.M{
			"resource": resource,
			"actions":  actions,
		})
	}

	return out
}

func privilegeSet() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cluster": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"db": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"collection": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"actions": &schema.Schema{
					Type:     schema.TypeSet,
					Required: true,
					MinItems: 1,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}
