package net_test

import (
	"fmt"
	"testing"
	"xeno/zohar/core/inet"
)

func Test_IPV4EndPoint_Functional_Basic(t *testing.T) {
	ipe := inet.NeoIPV4EndPointByStrIP(inet.EP_PROTO_TCP, 0, 255, "192.168.0.100", 10000)

	if !ipe.Valid() {
		t.Errorf("Invalid data of ipe %d", ipe.Identifier())
	}

	if ipe.IPV4Str() != "192.168.0.100" {
		t.Errorf("IP not match %s vs %s", ipe.IPV4Str(), "192.168.0.100")
	}
	ipe.SetIPV4Str("0.0.0.0")
	if ipe.IPV4Str() != "0.0.0.0" {
		t.Errorf("IP not match %s vs %s", ipe.IPV4Str(), "0.0.0.0")
	}

	if !ipe.Valid() {
		t.Errorf("Invalid data of ipe %d", ipe.Identifier())
	}

	ipe.SetIPV4Str("255.255.255.255")
	if ipe.IPV4Str() != "255.255.255.255" {
		t.Errorf("IP not match %s vs %s", ipe.IPV4Str(), "255.255.255.255")
	}

	if !ipe.Valid() {
		t.Errorf("Invalid data of ipe %d", ipe.Identifier())
	}

	if ipe.Port() != 10000 {
		t.Errorf("Port not match %d vs %d", ipe.Port(), 10000)
	}
	ipe.SetPort(0)
	if ipe.Port() != 0 {
		t.Errorf("Port not match %d vs %d", ipe.Port(), 0)
	}

	if !ipe.Valid() {
		t.Errorf("Invalid data of ipe %d", ipe.Identifier())
	}
	ipe.SetPort(65535)
	if ipe.Port() != 65535 {
		t.Errorf("Port not match %d vs %d", ipe.Port(), 65535)
	}

	if !ipe.Valid() {
		t.Errorf("Invalid data of ipe %d", ipe.Identifier())
	}

	if ipe.Extra() != 255 {
		t.Errorf("Extra not match %d vs %d", ipe.Extra(), 255)
	}
	ipe.SetExtra(0)
	if ipe.Extra() != 0 {
		t.Errorf("Extra not match %d vs %d", ipe.Extra(), 0)
	}

	if !ipe.Valid() {
		t.Errorf("Invalid data of ipe %d", ipe.Identifier())
	}

	if ipe.Mask() != 0 {
		t.Errorf("MASK not match %d vs %d", ipe.Mask(), 0)
	}
	ipe.SetMask(31)
	if ipe.Mask() != 31 {
		t.Errorf("MASK not match %d vs %d", ipe.Mask(), 63)
	}

	if !ipe.Valid() {
		t.Errorf("Invalid data of ipe %d", ipe.Identifier())
	}

	if ipe.Proto() != 1 {
		t.Errorf("Proto not match %d vs %d", ipe.Proto(), 0)
	}
	ipe.SetProto(2)
	if ipe.Extra() != 0 {
		t.Errorf("Proto not match %d vs %d", ipe.Proto(), 2)
	}

	if !ipe.Valid() {
		t.Errorf("Invalid data of ipe %d", ipe.Identifier())
	}

	fmt.Println(ipe.String())

}
