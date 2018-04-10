package gosoap

import (
	"testing"
)

var (
	tests = []struct {
		Params Params
		Err    string
	}{
		{
			Params: Params{"": ""},
			Err:    "error expected: xml: start tag with no name",
		},
	}
)

func TestClient_MarshalXML(t *testing.T) {
	soap, err := SoapClient("http://soapclient.com/xml/SQLDataSoap.WSDL")
	if err != nil {
		t.Errorf("error not expected: %s", err)
	}

	for _, test := range tests {
		_, err := soap.Call("SQLDataSRL", test.Params)
		if err == nil {
			t.Errorf(test.Err)
		}
	}
}
