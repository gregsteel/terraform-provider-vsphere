package vsphere

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/dvportgroup"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi/object"
)

func dataSourceVSphereDistributedPortGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereDistributedPortGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name or path of the distributed port group.",
				Optional:    true,
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "The managed object ID of the datacenter the distributed port group is in. This is not required when using ESXi directly, or if there is only one datacenter in your infrastructure.",
				Optional:    true,
			},
		},
	}
}

func dataSourceVSphereDistributedPortGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient

	name := d.Get("name").(string)
	if err := viapi.ValidateVirtualCenter(client); err == nil {
		if name == "" {
			return fmt.Errorf("name cannot be empty when using vCenter")
		}
	}

	var dc *object.Datacenter
	if dcID, ok := d.GetOk("datacenter_id"); ok {
		var err error
		dc, err = datacenterFromID(client, dcID.(string))
		if err != nil {
			return fmt.Errorf("cannot locate datacenter: %s", err)
		}
	}
	rp, err := dvportgroup.FromPartialName(client, name, dc)
	if err != nil {
		return fmt.Errorf("error fetching distributed port group: %s", err)
	}

	d.SetId(rp.Reference().Value)
	return nil
}
