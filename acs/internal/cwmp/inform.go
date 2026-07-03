package cwmp

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// ParseInform unmarshals a SOAP request body into an Inform plus the cwmp:ID header value.
func ParseInform(body []byte) (*Inform, string, error) {
	var env Envelope
	if err := xml.Unmarshal(body, &env); err != nil {
		return nil, "", fmt.Errorf("unmarshal soap envelope: %w", err)
	}
	if env.Body.Inform.DeviceID.SerialNumber == "" {
		return nil, "", fmt.Errorf("no Inform (or missing SerialNumber) in request body")
	}
	return &env.Body.Inform, env.Header.ID, nil
}

// Param returns the value of the named parameter from the ParameterList, or "" if absent.
func (in *Inform) Param(name string) string {
	for _, p := range in.ParameterList {
		if p.Name == name {
			return p.Value
		}
	}
	return ""
}

// ParamSuffix returns the first parameter value whose name ends with the given suffix.
// Useful because the Device model root differs across CPEs ("Device." vs "InternetGatewayDevice.").
func (in *Inform) ParamSuffix(suffix string) string {
	for _, p := range in.ParameterList {
		if strings.HasSuffix(p.Name, suffix) {
			return p.Value
		}
	}
	return ""
}

// FirstEventCode returns the EventCode of the first event, or "" if none.
func (in *Inform) FirstEventCode() string {
	if len(in.Events) == 0 {
		return ""
	}
	return in.Events[0].EventCode
}

// informResponseTmpl is a well-formed cwmp:InformResponse. %s is the echoed cwmp:ID.
const informResponseTmpl = `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
  <soap:Header>
    <cwmp:ID soap:mustUnderstand="1">%s</cwmp:ID>
  </soap:Header>
  <soap:Body>
    <cwmp:InformResponse>
      <MaxEnvelopes>1</MaxEnvelopes>
    </cwmp:InformResponse>
  </soap:Body>
</soap:Envelope>`

// BuildInformResponse returns a SOAP InformResponse echoing the given cwmp:ID.
func BuildInformResponse(cwmpID string) []byte {
	return fmt.Appendf(nil, informResponseTmpl, xmlEscape(cwmpID))
}

func xmlEscape(s string) string {
	var b strings.Builder
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}
