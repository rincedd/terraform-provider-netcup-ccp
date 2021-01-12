package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const HostURL string = "https://ccp.netcup.net/run/webservice/servers/endpoint.php?JSON"

type (
	CCPClient struct {
		HostURL    string
		AuthData   AuthData
		UserAgent  string
		httpClient http.Client
	}

	AuthData struct {
		CustomerNumber string `json:"customernumber"`
		APIKey         string `json:"apikey"`
		SessionId      string `json:"apisessionid"`
	}

	LoginData struct {
		CustomerNumber string `json:"customernumber"`
		APIKey         string `json:"apikey"`
		APIPassword    string `json:"apipassword"`
	}

	RequestBody struct {
		Action string      `json:"action"`
		Param  interface{} `json:"param"`
	}

	ResponseBody struct {
		ServerRequestId string `json:"serverrequestid"`
		Action          string `json:"action"`
		Status          string `json:"status"`     // Status of the Message like "error", "started", "pending", "warning" or "success".
		StatusCode      int    `json:"statuscode"` // Status code of the Message like 2011.
		ShortMessage    string `json:"shortmessage"`
		LongMessage     string `json:"longmessage"`
	}

	SessionData struct {
		SessionId string `json:"apisessionid"`
	}

	LoginResponse struct {
		ResponseBody
		ResponseData SessionData `json:"responsedata"`
	}

	DnsZone struct {
		Name         string `json:"domainname"`
		TTL          string `json:"ttl"`
		Serial       string `json:"serial"`
		Refresh      string `json:"refresh"`
		Retry        string `json:"retry"`
		Expire       string `json:"expire"`
		DNSSecStatus bool   `json:"dnssecstatus"`
	}

	DomainInfoRequest struct {
		AuthData
		DomainName string `json:"domainname"`
	}

	DnsZoneResponse struct {
		ResponseBody
		ResponseData DnsZone `json:"responsedata"`
	}

	DnsRecord struct {
		Id           string `json:"id,omitempty"`
		Hostname     string `json:"hostname"`
		Type         string `json:"type"`
		Priority     string `json:"priority,omitempty"`
		Destination  string `json:"destination"`
		DeleteRecord bool   `json:"deleterecord,omitempty"`
		State        string `json:"state,omitempty"`
		TTL          int    `json:"ttl,omitempty"`
	}

	DnsRecordsResponse struct {
		ResponseBody
		ResponseData struct {
			DnsRecords []DnsRecord `json:"dnsrecords,omitempty"`
		} `json:"responsedata"`
	}
)

func NewCCPClient(customerNumber, apiKey, apiPassword string) (*CCPClient, error) {
	c := CCPClient{
		HostURL:    HostURL,
		httpClient: http.Client{Timeout: 10 * time.Second},
	}

	err := c.login(customerNumber, apiKey, apiPassword)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *CCPClient) login(customerNumber, apiKey, apiPassword string) error {
	body, err := c.doRequest("login", LoginData{
		CustomerNumber: customerNumber,
		APIKey:         apiKey,
		APIPassword:    apiPassword,
	})
	res := LoginResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}

	c.AuthData = AuthData{
		CustomerNumber: customerNumber,
		APIKey:         apiKey,
		SessionId:      res.ResponseData.SessionId,
	}
	return nil
}

func (c *CCPClient) doRequest(action string, param interface{}) ([]byte, error) {
	rb, err := json.Marshal(RequestBody{
		Action: action,
		Param:  param,
	})

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.HostURL, strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (c *CCPClient) GetDnsZone(domainName string) (*DnsZone, error) {
	body, err := c.doRequest("infoDnsZone", DomainInfoRequest{
		AuthData:   c.AuthData,
		DomainName: domainName,
	})

	if err != nil {
		return nil, err
	}

	res := DnsZoneResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return &res.ResponseData, nil
}

func (c *CCPClient) GetDnsRecords(domainName string) ([]DnsRecord, error) {
	body, err := c.doRequest("infoDnsRecords", DomainInfoRequest{
		AuthData:   c.AuthData,
		DomainName: domainName,
	})

	if err != nil {
		return nil, err
	}

	res := DnsRecordsResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return res.ResponseData.DnsRecords, nil
}
