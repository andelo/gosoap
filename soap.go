package gosoap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Params type is used to set the params in soap request
type Params map[string]string

// SoapClient return new *Client to handle the requests with the WSDL
func SoapClient(wsdl string) (*Client, error) {
	_, err := url.Parse(wsdl)
	if err != nil {
		return nil, err
	}

	d, err := getWsdlDefinitions(wsdl)
	if err != nil {
		return nil, err
	}

	c := &Client{
		WSDL:        wsdl,
		URL:         strings.TrimSuffix(d.TargetNamespace, "/"),
		Definitions: d,
	}

	return c, nil
}

// Client struct hold all the informations about WSDL,
// request and response of the server
type Client struct {
	WSDL        string
	URL         string
	Method      string
	Params      Params
	Definitions *wsdlDefinitions
	Header      []interface{}
	payload     []byte
}

type Response struct {
	Body   []byte
	Header []byte
}

// Call call's the method m with Params p and optional soap headers h
func (c *Client) Call(m string, p Params, h ...interface{}) (res *Response, err error) {
	var payload []byte
	c.Method = m
	c.Params = p
	if len(h) > 0 {
		c.Header = h
	}

	payload, err = xml.MarshalIndent(c, "", "")
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.WriteString(xml.Header)
	buf.Write(payload)
	c.payload = buf.Bytes()

	b, err := c.doRequest()
	if err != nil {
		return nil, err
	}

	var soap SoapEnvelope
	err = xml.Unmarshal(b, &soap)

	res = &Response{
		Body:   soap.Body.Contents,
		Header: soap.Header.Contents,
	}

	return res, err
}

// Unmarshal get the body and unmarshal into the interface
func (r *Response) Unmarshal(v interface{}) error {
	if r == nil {
		return fmt.Errorf("response is nil")
	}

	if len(r.Body) == 0 {
		return fmt.Errorf("response body is empty")
	}

	var f Fault
	xml.Unmarshal(r.Body, &f)
	if f.Code != "" {
		return fmt.Errorf("[%s]: %s", f.Code, f.Description)
	}

	return xml.Unmarshal(r.Body, v)
}

// doRequest makes new request to the server using the c.Method, c.URL and the body.
// body is enveloped in Call method
func (c *Client) doRequest() ([]byte, error) {
	req, err := http.NewRequest("POST", c.WSDL, bytes.NewBuffer(c.payload))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	req.ContentLength = int64(len(c.payload))

	req.Header.Add("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Add("Accept", "text/xml")
	req.Header.Add("SOAPAction", fmt.Sprintf("%s/%s", c.URL, c.Method))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// SoapEnvelope struct
type SoapEnvelope struct {
	XMLName struct{} `xml:"Envelope"`
	Header  SoapHeader
	Body    SoapBody
}

// SoapHeader struct
type SoapHeader struct {
	XMLName  struct{} `xml:"Header"`
	Contents []byte   `xml:",innerxml"`
}

// SoapBody struct
type SoapBody struct {
	XMLName  struct{} `xml:"Body"`
	Contents []byte   `xml:",innerxml"`
}
