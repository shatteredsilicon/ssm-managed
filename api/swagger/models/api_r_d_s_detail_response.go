// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// APIRDSDetailResponse api r d s detail response
// swagger:model apiRDSDetailResponse
type APIRDSDetailResponse struct {

	// aws access key id
	AwsAccessKeyID string `json:"aws_access_key_id,omitempty"`

	// aws secret access key
	AwsSecretAccessKey string `json:"aws_secret_access_key,omitempty"`

	// instance
	Instance string `json:"instance,omitempty"`

	// region
	Region string `json:"region,omitempty"`
}

// Validate validates this api r d s detail response
func (m *APIRDSDetailResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *APIRDSDetailResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *APIRDSDetailResponse) UnmarshalBinary(b []byte) error {
	var res APIRDSDetailResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
