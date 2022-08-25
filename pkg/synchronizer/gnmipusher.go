// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

// GNMIPusher implements a pusher that pushes to a REST API endpoint.

package synchronizer

import (
	"context"
	"fmt"
	gnmiapi "github.com/openconfig/gnmi/proto/gnmi"
	"os"
)

// GNMIPusher implements a pusher that pushes to a gnmi endpoint.
type GNMIPusher struct {
	endpoint   string
	payload    string
	target     string
	pushClient Client
}

const (
	// StringVal :
	StringVal = "string_val"

	// IntVal :
	IntVal = "int_val"

	// BoolVal :
	BoolVal = "bool_val"
)

// NewGNMIPusher allocates a gnmi pusher for a given endpoint
func NewGNMIPusher(url string, target string, payload string) PusherInterface {
	return NewGNMIPusherWithClient(url, payload, target, &client{})
}

// NewGNMIPusherWithClient allocates a gnmi pusher for a given endpoint
func NewGNMIPusherWithClient(url string, target string, payload string, pushClient Client) PusherInterface {
	gnmiPusher := &GNMIPusher{
		endpoint:   url,
		pushClient: pushClient,
		payload:    payload,
		target:     target,
	}

	return gnmiPusher
}

// PushUpdate pushes an update to the GNMI server.
func (p *GNMIPusher) PushUpdate() error {
	setGnmiRequest := &gnmiapi.SetRequest{}

	// update:{path:{elem:{name:"someURL"} target:"stratum"} val:{string_val:"somepayload"}}
	e := &gnmiapi.PathElem{
		Name: p.endpoint,
	}
	es := []*gnmiapi.PathElem{e}
	path := &gnmiapi.Path{
		Origin: "",
		Elem:   es,
		Target: p.target,
	}
	tv := &gnmiapi.TypedValue{
		Value: &gnmiapi.TypedValue_StringVal{
			StringVal: p.payload,
		},
	}
	ud := &gnmiapi.Update{
		Path:       path,
		Val:        tv,
		Duplicates: 0,
	}
	uds := []*gnmiapi.Update{ud}
	//var protoBuilder strings.Builder
	//protoBuilder.WriteString("update:{path:{elem:{name:\"" + p.endpoint + "\"}")
	//protoBuilder.WriteString("  target:\"" + p.target + "\"}")
	//protoBuilder.WriteString("val:{string_val:\"" + p.payload + "\"}}")
	//protoString := protoBuilder.String()

	setGnmiRequest.Update = uds
	//if err := proto.UnmarshalText(protoString, setGnmiRequest); err != nil {
	//	return err
	//}

	fmt.Fprintf(os.Stderr, "proto for setGnmiRequest is: %v", setGnmiRequest)

	_, err := p.pushClient.Set(context.Background(), setGnmiRequest)
	if err != nil {
		return &PushError{
			Endpoint:   p.endpoint,
			StatusCode: 500, // Not sure what the right thing to do is
			Status:     err.Error(),
			Operation:  "SET",
		}
	}
	return nil
}

// PushDelete pushes a delete operation to the GNMI server
func (p *GNMIPusher) PushDelete() error {
	return nil
}