// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
	"fmt"
)

// just to fix compile issues with the import
var (
	_ fmt.Formatter
	_ context.Context
)

type DynamicDNS struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	CustomService string   `json:"custom_service,omitempty"` // ^[^"' ]+$
	HostName      string   `json:"host_name,omitempty"`      // ^[^"' ]+$
	Interface     string   `json:"interface,omitempty"`      // wan|wan2
	Login         string   `json:"login,omitempty"`          // ^[^"' ]+$
	Options       []string `json:"options,omitempty"`        // ^[^"' ]+$
	Server        string   `json:"server"`                   // ^[^"' ]+$|^$
	Service       string   `json:"service,omitempty"`        // afraid|changeip|cloudflare|dnspark|dslreports|dyndns|easydns|googledomains|namecheap|noip|sitelutions|zoneedit|custom
	XPassword     string   `json:"x_password,omitempty"`     // ^[^"' ]+$
}

func (c *Client) listDynamicDNS(ctx context.Context, site string) ([]DynamicDNS, error) {
	var respBody struct {
		Meta meta         `json:"meta"`
		Data []DynamicDNS `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("proxy/network/api/s/%s/rest/dynamicdns", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getDynamicDNS(ctx context.Context, site, id string) (*DynamicDNS, error) {
	var respBody struct {
		Meta meta         `json:"meta"`
		Data []DynamicDNS `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("proxy/network/api/s/%s/rest/dynamicdns/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteDynamicDNS(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("proxy/network/api/s/%s/rest/dynamicdns/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createDynamicDNS(ctx context.Context, site string, d *DynamicDNS) (*DynamicDNS, error) {
	var respBody struct {
		Meta meta         `json:"meta"`
		Data []DynamicDNS `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("proxy/network/api/s/%s/rest/dynamicdns", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateDynamicDNS(ctx context.Context, site string, d *DynamicDNS) (*DynamicDNS, error) {
	var respBody struct {
		Meta meta         `json:"meta"`
		Data []DynamicDNS `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("proxy/network/api/s/%s/rest/dynamicdns/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}
