package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rincedd/terraform-provider-netcup-ccp/internal/client"
)

func dataSourceDnsZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDnsZoneRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"serial": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"refresh": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retry": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_sec_status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceDnsZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	domainName := d.Get("name").(string)
	ccpClient := m.(*client.CCPClient)

	if ccpClient == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Missing Netcup CCP client",
			Detail:   "Unable to retrieve DNS zone data without Netcup CCP client",
		})
		return diags
	}

	dnsZone, err := ccpClient.GetDnsZone(domainName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to retrieve data for DNS zone",
			Detail:   fmt.Sprintf("Unable to retrieve data for DNS zone %s [%s]", domainName, err.Error()),
		})
		return diags
	}

	d.SetId(domainName)
	d.Set("ttl", dnsZone.TTL)
	d.Set("serial", dnsZone.Serial)
	d.Set("refresh", dnsZone.Refresh)
	d.Set("retry", dnsZone.Retry)
	d.Set("expire", dnsZone.Expire)
	d.Set("dns_sec_status", dnsZone.DNSSecStatus)

	return diags
}
