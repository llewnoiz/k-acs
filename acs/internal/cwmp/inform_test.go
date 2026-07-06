package cwmp

import (
	"os"
	"strings"
	"testing"
)

func TestParseInform(t *testing.T) {
	body, err := os.ReadFile("../../testdata/inform_sample.xml")
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}

	in, cwmpID, err := ParseInform(body)
	if err != nil {
		t.Fatalf("ParseInform: %v", err)
	}

	if cwmpID != "100" {
		t.Errorf("cwmpID = %q, want %q", cwmpID, "100")
	}
	if in.DeviceID.SerialNumber != "SN-0001" {
		t.Errorf("SerialNumber = %q, want %q", in.DeviceID.SerialNumber, "SN-0001")
	}
	if in.DeviceID.Manufacturer != "Acme Networks" {
		t.Errorf("Manufacturer = %q", in.DeviceID.Manufacturer)
	}
	if in.DeviceID.ProductClass != "Router-X100" {
		t.Errorf("ProductClass = %q", in.DeviceID.ProductClass)
	}
	if got := in.FirstEventCode(); got != "2 PERIODIC" {
		t.Errorf("FirstEventCode = %q", got)
	}
	if got := in.ParamSuffix("DeviceInfo.SoftwareVersion"); got != "1.4.2" {
		t.Errorf("SoftwareVersion = %q, want 1.4.2", got)
	}
	if got := in.ParamSuffix("DeviceInfo.HardwareVersion"); got != "revB" {
		t.Errorf("HardwareVersion = %q, want revB", got)
	}
	if got := in.ParamSuffix("ManagementServer.ConnectionRequestURL"); !strings.Contains(got, "192.168.0.10") {
		t.Errorf("ConnectionRequestURL = %q", got)
	}

	if got := in.ParamSuffix("DeviceInfo.UpTime"); got != "123456" {
		t.Errorf("UpTime = %q, want 123456", got)
	}

	if n := len(in.ParameterList); n != 4 {
		t.Errorf("ParameterList len = %d, want 4", n)
	}
}

func TestBuildInformResponse(t *testing.T) {
	out := string(BuildInformResponse("100"))
	if !strings.Contains(out, "cwmp:InformResponse") {
		t.Errorf("response missing InformResponse element:\n%s", out)
	}
	if !strings.Contains(out, ">100<") {
		t.Errorf("response missing echoed cwmp:ID:\n%s", out)
	}
	if !strings.Contains(out, "<MaxEnvelopes>1</MaxEnvelopes>") {
		t.Errorf("response missing MaxEnvelopes:\n%s", out)
	}
}
