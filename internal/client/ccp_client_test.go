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

		So(client.authData.SessionId, ShouldEqual, "SESSION_ID")
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

func TestCCPClient_CreateDnsRecord(t *testing.T) {
	Convey("creates new DNS record and returns it", t, func() {
		client, tearDown := setupClientTest()
		defer tearDown()

		gock.New(HostURL).Post("").
			BodyString(`{"action":"updateDnsRecords","param":{"customernumber":"CUSTOMER_NUMBER","apikey":"API_KEY","apisessionid":"SESSION_ID","domainname":"domain.com","dnsrecordset":{"dnsrecords":[{"hostname":"HOSTNAME","type":"TXT","destination":"DESTINATION"}]}}}`).
			Reply(200).Type("application/json").
			BodyString(`{"responsedata":{"dnsrecords":[{"id":"5838738","hostname":"*","type":"A","priority":"0","destination":"1.2.3.4","deleterecord":false,"state":"yes"},{"id":"5838739","hostname":"HOSTNAME","type":"TXT","priority":"0","destination":"DESTINATION","deleterecord":false,"state":"yes"}]}}`)

		newRecord, err := client.CreateDnsRecord("domain.com", NewDnsRecord{
			Hostname:    "HOSTNAME",
			Type:        "TXT",
			Destination: "DESTINATION",
		})

		So(err, ShouldBeNil)
		So(*newRecord, ShouldResemble, DnsRecord{
			Id:           "5838739",
			Hostname:     "HOSTNAME",
			Type:         "TXT",
			Priority:     "0",
			Destination:  "DESTINATION",
			DeleteRecord: false,
			State:        "yes",
		})
	})
}
