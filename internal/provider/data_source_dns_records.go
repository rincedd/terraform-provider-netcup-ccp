package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rincedd/terraform-provider-netcup-ccp/internal/client"
)

func dataSourceDnsRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDnsRecordsRead,
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDnsRecordsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	domainName := d.Get("domain_name").(string)
	ccpClient := m.(*client.CCPClient)

	if ccpClient == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Missing Netcup CCP client",
			Detail:   "Unable to retrieve DNS records without Netcup CCP client",
		})
		return diags
	}

	dnsRecords, err := ccpClient.GetDnsRecords(domainName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to retrieve DNS records",
			Detail:   fmt.Sprintf("Unable to retrieve DNS records for %s: %s", domainName, err.Error()),
		})
		return diags
	}

	if dnsRecords != nil {
		records := make([]interface{}, len(dnsRecords), len(dnsRecords))

		for i, dnsRecord := range dnsRecords {
			record := make(map[string]interface{})
			record["id"] = dnsRecord.Id
			record["name"] = dnsRecord.Hostname
			record["type"] = dnsRecord.Type
			record["value"] = dnsRecord.Destination
			record["priority"] = dnsRecord.Priority
			record["state"] = dnsRecord.State
			records[i] = record
		}

		d.SetId(domainName)
		if err := d.Set("records", records); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
