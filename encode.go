package gosoap

import (
	"encoding/xml"
	"fmt"
	"reflect"
)

var tokens []xml.Token

// MarshalXML envelope the body and encode to xml
func (c Client) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {

	tokens = []xml.Token{}

	//start envelope
	if c.Definitions == nil {
		return fmt.Errorf("definitions is nil")
	}

	startEnvelope()
	if c.Header != nil && len(c.Header) > 0 {
		startHeader()
		for _, h := range c.Header {
			header := reflect.ValueOf(h)
			t := xml.StartElement{
				Name: xml.Name{
					Space: c.Definitions.TargetNamespace,
					Local: reflect.Indirect(header).Type().Name(),
				},
			}
			tokens = append(tokens, t)
			for i := 0; i < header.NumField(); i++ {
				name := header.Type().Field(i).Name
				value := fmt.Sprintf("%v", header.Field(i))

				t := xml.StartElement{
					Name: xml.Name{
						Space: "",
						Local: name,
					},
				}
				tokens = append(tokens, t, xml.CharData(value), xml.EndElement{Name: t.Name})
			}
			tokens = append(tokens, xml.EndElement{Name: t.Name})
		}
		endHeader()
	}

	err := startBody(c.Method, c.Definitions.TargetNamespace)
	if err != nil {
		return err
	}

	for k, v := range c.Params {
		if k != "" {
			t := xml.StartElement{
				Name: xml.Name{
					Space: "",
					Local: k,
				},
			}

			tokens = append(tokens, t, xml.CharData(v), xml.EndElement{Name: t.Name})
		}
	}
	//end envelope
	endBody(c.Method)
	endEnvelope()

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	return e.Flush()
}

func startEnvelope() {
	e := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "soap:Envelope",
		},
		Attr: []xml.Attr{
			{Name: xml.Name{Space: "", Local: "xmlns:xsi"}, Value: "http://www.w3.org/2001/XMLSchema-instance"},
			{Name: xml.Name{Space: "", Local: "xmlns:xsd"}, Value: "http://www.w3.org/2001/XMLSchema"},
			{Name: xml.Name{Space: "", Local: "xmlns:soap"}, Value: "http://schemas.xmlsoap.org/soap/envelope/"},
		},
	}

	tokens = append(tokens, e)
}

func endEnvelope() {
	e := xml.EndElement{
		Name: xml.Name{
			Space: "",
			Local: "soap:Envelope",
		},
	}

	tokens = append(tokens, e)
}

func startHeader() {
	h := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "soap:Header",
		},
	}

	tokens = append(tokens, h)
}

func endHeader() {
	h := xml.EndElement{
		Name: xml.Name{
			Space: "",
			Local: "soap:Header",
		},
	}

	tokens = append(tokens, h)
}

// startToken initiate body of the envelope
func startBody(m, n string) error {
	b := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "soap:Body",
		},
	}

	if m == "" || n == "" {
		return fmt.Errorf("method or namespace is empty")
	}

	r := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: m,
		},
		Attr: []xml.Attr{
			{Name: xml.Name{Space: "", Local: "xmlns"}, Value: n},
		},
	}

	tokens = append(tokens, b, r)

	return nil
}

// endToken close body of the envelope
func endBody(m string) {
	b := xml.EndElement{
		Name: xml.Name{
			Space: "",
			Local: "soap:Body",
		},
	}

	r := xml.EndElement{
		Name: xml.Name{
			Space: "",
			Local: m,
		},
	}

	tokens = append(tokens, r, b)
}
