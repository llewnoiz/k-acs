// Package cwmp implements minimal TR-069 (CWMP) SOAP handling for the Inform RPC.
package cwmp

import "encoding/xml"

// Envelope is the SOAP envelope wrapping a CWMP message.
// We only model the pieces needed to read an Inform and to know the cwmp:ID header.
type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Header  Header   `xml:"Header"`
	Body    Body     `xml:"Body"`
}

// Header carries the CWMP message ID (echoed back in the response).
type Header struct {
	ID string `xml:"ID"`
}

// Body holds the Inform RPC (the only inbound RPC we handle in this POC).
type Body struct {
	Inform Inform `xml:"Inform"`
}

// Inform is the cwmp:Inform RPC a CPE sends when it starts a session.
type Inform struct {
	DeviceID      DeviceID               `xml:"DeviceId"`
	Events        []EventStruct          `xml:"Event>EventStruct"`
	MaxEnvelopes  int                    `xml:"MaxEnvelopes"`
	CurrentTime   string                 `xml:"CurrentTime"`
	RetryCount    int                    `xml:"RetryCount"`
	ParameterList []ParameterValueStruct `xml:"ParameterList>ParameterValueStruct"`
}

// DeviceID uniquely identifies a CPE.
type DeviceID struct {
	Manufacturer string `xml:"Manufacturer"`
	OUI          string `xml:"OUI"`
	ProductClass string `xml:"ProductClass"`
	SerialNumber string `xml:"SerialNumber"`
}

// EventStruct describes why the CPE contacted the ACS (e.g. "0 BOOTSTRAP", "2 PERIODIC").
type EventStruct struct {
	EventCode  string `xml:"EventCode"`
	CommandKey string `xml:"CommandKey"`
}

// ParameterValueStruct is a single name/value pair from the Inform ParameterList.
type ParameterValueStruct struct {
	Name  string `xml:"Name"`
	Value string `xml:"Value"`
}
