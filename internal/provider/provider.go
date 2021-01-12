package provider

import (
	"context"
	"github.com/rincedd/terraform-provider-netcup-ccp/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"customer_number": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("NETCUP_CUSTOMER_NUMBER", nil),
					Description: "Netcup customer number.",
				},
				"ccp_api_key": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("NETCUP_CCP_API_KEY", nil),
					Description: "Netcup CCP API key.",
				},
				"ccp_api_password": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("NETCUP_CCP_API_PASSWORD", nil),
					Description: "Netcup CCP API password.",
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"netcup_ccp_dns_zone":    dataSourceDnsZone(),
				"netcup_ccp_dns_records": dataSourceDnsRecords(),
			},
			ResourcesMap: map[string]*schema.Resource{},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		customerNumber := d.Get("customer_number").(string)
		ccpApiKey := d.Get("ccp_api_key").(string)
		ccpApiPassword := d.Get("ccp_api_password").(string)

		var diags diag.Diagnostics

		if (customerNumber != "") && (ccpApiKey != "") && (ccpApiPassword != "") {
			ccpClient, err := client.NewCCPClient(customerNumber, ccpApiKey, ccpApiPassword)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Unable to authenticate with CCP API",
					Detail:   "Unable to authenticate customer " + customerNumber + " with Netcup CCP API",
				})
				return nil, diags
			}

			userAgent := p.UserAgent("terraform-provider-netcup-ccp", version)
			ccpClient.UserAgent = userAgent

			return ccpClient, nil
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Missing customer number/API key/API password.",
			Detail:   "Netcup customer number, API key, and API password are required.",
		})
		return nil, diags
	}
}
