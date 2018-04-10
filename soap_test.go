package gosoap

import (
	"fmt"
	"testing"
)

var (
	scts = []struct {
		URL string
		Err bool
	}{
		{
			URL: "://www.server",
			Err: false,
		},
		{
			URL: "",
			Err: false,
		},
		{
			URL: "http://soapclient.com/xml/SQLDataSoap.WSDL",
			Err: true,
		},
	}
)

func TestSoapClient(t *testing.T) {
	for _, sct := range scts {
		_, err := SoapClient(sct.URL)
		if err != nil && sct.Err {
			t.Errorf("URL: %s - error: %s", sct.URL, err)
		}
	}
}

type SQLDataSRLResponse struct {
	SQLDataSRLResult SQLDataSRLResult
}

type SQLDataSRLResult struct {
	ReturnCode        string
	IP                string
	ReturnCodeDetails string
	CountryName       string
	CountryCode       string
}

var (
	r SQLDataSRLResponse

	params = Params{"SRLFile": "WHOIS.SRI", "RequestName": "whois"}
)

func TestClient_Call(t *testing.T) {
	var (
		soap *Client
		res  *Response
		err  error
	)
	soap, err = SoapClient("http://soapclient.com/xml/SQLDataSoap.WSDL")
	if err != nil {
		t.Errorf("error not expected: %s", err)
	}
	soap.URL = ""
	params["key"] = "apple.com"
	res, err = soap.Call("", params)
	if err == nil {
		t.Errorf("method is empty")
	}

	err = res.Unmarshal(&r)
	if err == nil {
		t.Errorf("body is empty")
	}

	res, err = soap.Call("SQLDataSRL", params)
	if err != nil {
		t.Errorf("error in soap call: %s", err)
	}

	fmt.Println(string(res.Body))
	res.Unmarshal(&r)
	if r.SQLDataSRLResult.CountryCode != "USA" {
		t.Errorf("error: %+v", r)
	}

	c := &Client{}
	_, err = c.Call("", Params{})
	if err == nil {
		t.Errorf("error expected but nothing got.")
	}

	c.WSDL = "://test."

	_, err = c.Call("SQLDataSRL", params)
	if err == nil {
		t.Errorf("invalid WSDL")
	}
}

func TestClient_doRequest(t *testing.T) {
	c := &Client{}

	_, err := c.doRequest()
	if err == nil {
		t.Errorf("body is empty")
	}

	c.WSDL = "://teste."
	_, err = c.doRequest()
	if err == nil {
		t.Errorf("invalid WSDL")
	}
}
