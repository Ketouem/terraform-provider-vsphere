package vsphere

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceVSphereNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereNetworkRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name or path of the network.",
				Required:    true,
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "The managed object ID of the datacenter the network is in. This is required if the supplied path is not an absolute path containing a datacenter and there are multiple datacenters in your infrastructure.",
				Optional:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The managed object type of the network.",
				Computed:    true,
			},
		},
	}
}

func dataSourceVSphereNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*VSphereClient).vimClient
	if err := validateVirtualCenter(client); err != nil {
		return err
	}

	name := d.Get("name").(string)
	dc, err := datacenterFromID(client, d.Get("datacenter_id").(string))
	if err != nil {
		return fmt.Errorf("cannot locate datacenter: %s", err)
	}
	net, err := networkFromPath(client, name, dc)
	if err != nil {
		return fmt.Errorf("error fetching network: %s", err)
	}

	d.SetId(net.Reference().Value)
	d.Set("type", net.Reference().Type)
	return nil
}
