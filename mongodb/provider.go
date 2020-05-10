package mongodb

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var mongo_uri string

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"mongo_uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A MongoDB connection string",
			},
		},
		ConfigureFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"mongodb_role": roleResourceServer(),
			"mongodb_user": userResourceServer(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	connectionString, _ := d.GetOk("mongo_uri") //dTos("mongo_uri", d)
	mongo_uri = d.Get("mongo_uri").(string)

	return NewClient(connectionString.(string))
}
