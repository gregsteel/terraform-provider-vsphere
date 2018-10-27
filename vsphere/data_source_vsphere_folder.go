package vsphere

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/folder"
)

func dataSourceVSphereFolder() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereFolderRead,

		Schema: map[string]*schema.Schema{
			"path": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The path of the folder.",
				Optional:    true,
			},
		},
	}
}

func dataSourceVSphereFolderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	path := d.Get("path").(string)
	hs, err := folder.FromAbsolutePath(client, path)
	if err != nil {
		return fmt.Errorf("error fetching folder: %s", err)
	}
	id := hs.Reference().Value
	d.SetId(id)
	return nil
}
