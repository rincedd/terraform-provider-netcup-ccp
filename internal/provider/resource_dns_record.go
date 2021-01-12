package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rincedd/terraform-provider-netcup-ccp/internal/client"
)

func resourceDnsRecord() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "DNS record",

		CreateContext: resourceDnsRecordCreate,
		ReadContext:   resourceDnsRecordRead,
		UpdateContext: resourceDnsRecordUpdate,
		DeleteContext: resourceDnsRecordDelete,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
		},
	}
}

func resourceDnsRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	record := client.NewDnsRecord{
		Hostname:    d.Get("name").(string),
		Type:        d.Get("type").(string),
		Destination: d.Get("value").(string),
		Priority:    d.Get("priority").(string),
	}
	domainName := d.Get("domain_name").(string)
	ccpClient := m.(*client.CCPClient)

	newRecord, err := ccpClient.CreateDnsRecord(domainName, record)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newRecord.Id)

	return nil
}

func resourceDnsRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domainName := d.Get("domain_name").(string)
	ccpClient := m.(*client.CCPClient)

	record, err := ccpClient.GetDnsRecordById(domainName, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", record.Hostname)
	d.Set("type", record.Type)
	d.Set("value", record.Destination)
	d.Set("priority", record.Priority)

	return nil
}

func resourceDnsRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domainName := d.Get("domain_name").(string)
	ccpClient := m.(*client.CCPClient)

	record, err := ccpClient.UpdateDnsRecord(domainName, client.DnsRecord{
		Id:           d.Id(),
		Hostname:     d.Get("name").(string),
		Type:         d.Get("type").(string),
		Priority:     d.Get("priority").(string),
		Destination:  d.Get("value").(string),
		DeleteRecord: false,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", record.Hostname)
	d.Set("type", record.Type)
	d.Set("value", record.Destination)
	d.Set("priority", record.Priority)

	return nil
}

func resourceDnsRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domainName := d.Get("domain_name").(string)
	ccpClient := m.(*client.CCPClient)

	err := ccpClient.DeleteDnsRecord(domainName, client.DnsRecord{
		Id:          d.Id(),
		Hostname:    d.Get("name").(string),
		Type:        d.Get("type").(string),
		Priority:    d.Get("priority").(string),
		Destination: d.Get("value").(string),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
