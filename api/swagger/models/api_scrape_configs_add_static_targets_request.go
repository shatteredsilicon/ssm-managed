// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// APIScrapeConfigsAddStaticTargetsRequest api scrape configs add static targets request
// swagger:model apiScrapeConfigsAddStaticTargetsRequest

type APIScrapeConfigsAddStaticTargetsRequest struct {

	// Check that added targets can be scraped from PMM Server
	CheckReachability bool `json:"check_reachability,omitempty"`

	// job name
	JobName string `json:"job_name,omitempty"`

	// Hostnames or IPs followed by an optional port number: "1.2.3.4:9090"
	Targets []string `json:"targets"`
}

/* polymorph apiScrapeConfigsAddStaticTargetsRequest check_reachability false */

/* polymorph apiScrapeConfigsAddStaticTargetsRequest job_name false */

/* polymorph apiScrapeConfigsAddStaticTargetsRequest targets false */

// Validate validates this api scrape configs add static targets request
func (m *APIScrapeConfigsAddStaticTargetsRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTargets(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *APIScrapeConfigsAddStaticTargetsRequest) validateTargets(formats strfmt.Registry) error {

	if swag.IsZero(m.Targets) { // not required
		return nil
	}

	return nil
}

// MarshalBinary interface implementation
func (m *APIScrapeConfigsAddStaticTargetsRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *APIScrapeConfigsAddStaticTargetsRequest) UnmarshalBinary(b []byte) error {
	var res APIScrapeConfigsAddStaticTargetsRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
