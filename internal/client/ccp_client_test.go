package client

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

const (
	customerNumber = "CUSTOMER_NUMBER"
	apiKey         = "API_KEY"
	apiPassword    = "API_PASSWORD"
)

func setupClientTest() (*CCPClient, func()) {
	gock.New(HostURL).Post("").BodyString(`{"action":"login","param":{"customernumber":"CUSTOMER_NUMBER","apikey":"API_KEY","apipassword":"API_PASSWORD"}}`).
		Reply(200).Type("application/json").
		BodyString(`{"responsedata":{"apisessionid":"SESSION_ID"}}`)

	client, err := NewCCPClient(customerNumber, apiKey, apiPassword)

	if err != nil {
		panic(err)
	}

	return client, gock.Off
}

func TestNewCCPClient(t *testing.T) {
	Convey("sends a login request upon creation", t, func() {
		client, tearDown := setupClientTest()
		defer tearDown()

		So(client.AuthData.SessionId, ShouldEqual, "SESSION_ID")
	})
}

func TestCCPClient_GetDnsZone(t *testing.T) {
	Convey("retrieves the DNS zone data for a given domain name", t, func() {
		client, tearDown := setupClientTest()
		defer tearDown()

		gock.New(HostURL).Post("").
			BodyString(`{"action":"infoDnsZone","param":{"customernumber":"CUSTOMER_NUMBER","apikey":"API_KEY","apisessionid":"SESSION_ID","domainname":"domain.com"}}`).
			Reply(200).Type("application/json").
			BodyString(`{"responsedata":{"domainname":"domain.com","ttl":"86400","serial":"1234","refresh":"28800","retry":"7200","expire":"1209600","dnssecstatus":true}}`)

		dnsZone, err := client.GetDnsZone("domain.com")

		So(err, ShouldBeNil)
		So(*dnsZone, ShouldResemble, DnsZone{
			Name:         "domain.com",
			TTL:          "86400",
			Serial:       "1234",
			Refresh:      "28800",
			Retry:        "7200",
			Expire:       "1209600",
			DNSSecStatus: true,
		})
	})
}
