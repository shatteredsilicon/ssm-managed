// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// InventoryAddPMMAgentResponse inventory add PMM agent response
// swagger:model inventoryAddPMMAgentResponse
type InventoryAddPMMAgentResponse struct {

	// pmm agent
	PMMAgent *InventoryPMMAgent `json:"pmm_agent,omitempty"`

	// pmm-agent UUID.
	UUID string `json:"uuid,omitempty"`
}

// Validate validates this inventory add PMM agent response
func (m *InventoryAddPMMAgentResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePMMAgent(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *InventoryAddPMMAgentResponse) validatePMMAgent(formats strfmt.Registry) error {

	if swag.IsZero(m.PMMAgent) { // not required
		return nil
	}

	if m.PMMAgent != nil {
		if err := m.PMMAgent.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("pmm_agent")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *InventoryAddPMMAgentResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *InventoryAddPMMAgentResponse) UnmarshalBinary(b []byte) error {
	var res InventoryAddPMMAgentResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
