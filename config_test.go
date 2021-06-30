package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/SophosNSG/paloma-config-service/src"
)

func TestParseAndPublishDNSConfig(t *testing.T) {
	reqBody := []byte(`{   
		"DNSFwdSelection": "5", 
		"dns1": "10.43.77.148",
		"dns2": "192.168.0.210",
		"dns3": "192.168.0.209",
		"ipv6dns1" : "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"ipv6dns2" : "2001:0db8:85a3:0000:0000:8a2e:0370:7335",
		"ipv6dns3" : "2001:0db8:85a3:0000:0000:8a2e:0370:7336",
		"ipv6rdodhcpserver" : "1",
		"rdodhcpserver": "3"
	}`) // request body "DNSFwdSelection": [0-3] Should throw blank response with error

	resp, err := src.ParseAndPublishDNSConfig(reqBody)

	if resp != "" || err == nil {
		t.Errorf("Response was incorrect or error was empty, should be empty")
	}

	var data map[string]interface{}
	err = json.Unmarshal(reqBody, &data)
	if err != nil {
		t.Errorf("Error occurred while unmarshaling")
	}

	data["DNSFwdSelection"] = "3"   // Reverting back the DNWFwdSelection to Normal
	data["dns1"] = "10.43.77.148.1" // Setting dns1 with invalid IP

	reqBody, _ = json.Marshal(data)

	fmt.Println(fmt.Sprintf("Marshalled json obj is %s", string(reqBody)))

	resp, err = src.ParseAndPublishDNSConfig(reqBody)

	if resp != "" || err == nil {
		t.Errorf("Response was incorrect or error was empty, should be empty")
	}

	data["dns1"] = "10.43.77.148"                                     // Reverting this to normal
	data["ipv6dns3"] = "2001:0db8:85a3:0000:0000:8a2e:0370:7335:1789" // Should throw an error if

	reqBody, _ = json.Marshal(data)

	fmt.Println(fmt.Sprintf("Marshalled json obj is %s", string(reqBody)))

	resp, err = src.ParseAndPublishDNSConfig(reqBody)

	if resp != "" || err == nil {
		t.Errorf("Response was incorrect or error was empty, should be empty")
	}

	data["ipv6dns3"] = "2001:0db8:85a3:0000:0000:8a2e:0370:7335"

	reqBody, _ = json.Marshal(data)

	fmt.Println(fmt.Sprintf("Marshalled json obj is %s", string(reqBody)))

	resp, err = src.ParseAndPublishDNSConfig(reqBody) // The publisher needs to be mocked before being called

	fmt.Println(fmt.Sprintf("Response got from the function is %s", resp))
	if resp != "" || err == nil {
		t.Errorf("Response was missing or error was present")
	}

}

func TestParseAndWriteDNSRecodConfig(t *testing.T) {
	reqBody := []byte(`{
		"fqdn": "sophosbuddy.com",
		"reverseresolve": true,
		"ip_list": [{
			"ipfamily": 0,
			"iptype": 0,
			"ipaddress": "10.194.194.51.1",
			"ttl": 59,
			"weight": 43,
			"isWan": false
		}]
	}`)

	fmt.Println(fmt.Sprintf("Marshalled json obj is %s", string(reqBody)))

	resp, err := src.ParseAndWriteDNSRecodConfig(reqBody)

	if resp != "" || err == nil {
		t.Errorf("Response was incorrect or error was empty, should be empty")
	}

	reqBody = []byte(`{
		"fqdn": "sophosbuddy.com",
		"reverseresolve": true,
		"ip_list": [{
			"ipfamily": 0,
			"iptype": 0,
			"ipaddress": "10.194.194.51",
			"ttl": 59,
			"weight": 43,
			"isWan": false
		}]
	}`)

	fmt.Println(fmt.Sprintf("Marshalled json obj is %s", string(reqBody)))

	resp, err = src.ParseAndWriteDNSRecodConfig(reqBody)

	if resp != "" || err == nil {
		t.Errorf("Response was incorrect or error was empty, should be empty")
	}
}
