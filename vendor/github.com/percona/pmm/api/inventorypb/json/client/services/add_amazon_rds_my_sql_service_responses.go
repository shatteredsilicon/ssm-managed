// Code generated by go-swagger; DO NOT EDIT.

package services

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// AddAmazonRDSMySQLServiceReader is a Reader for the AddAmazonRDSMySQLService structure.
type AddAmazonRDSMySQLServiceReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AddAmazonRDSMySQLServiceReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewAddAmazonRDSMySQLServiceOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewAddAmazonRDSMySQLServiceDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewAddAmazonRDSMySQLServiceOK creates a AddAmazonRDSMySQLServiceOK with default headers values
func NewAddAmazonRDSMySQLServiceOK() *AddAmazonRDSMySQLServiceOK {
	return &AddAmazonRDSMySQLServiceOK{}
}

/*AddAmazonRDSMySQLServiceOK handles this case with default header values.

A successful response.
*/
type AddAmazonRDSMySQLServiceOK struct {
	Payload *AddAmazonRDSMySQLServiceOKBody
}

func (o *AddAmazonRDSMySQLServiceOK) Error() string {
	return fmt.Sprintf("[POST /v1/inventory/Services/AddAmazonRDSMySQL][%d] addAmazonRdsMySqlServiceOk  %+v", 200, o.Payload)
}

func (o *AddAmazonRDSMySQLServiceOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(AddAmazonRDSMySQLServiceOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAddAmazonRDSMySQLServiceDefault creates a AddAmazonRDSMySQLServiceDefault with default headers values
func NewAddAmazonRDSMySQLServiceDefault(code int) *AddAmazonRDSMySQLServiceDefault {
	return &AddAmazonRDSMySQLServiceDefault{
		_statusCode: code,
	}
}

/*AddAmazonRDSMySQLServiceDefault handles this case with default header values.

An error response.
*/
type AddAmazonRDSMySQLServiceDefault struct {
	_statusCode int

	Payload *AddAmazonRDSMySQLServiceDefaultBody
}

// Code gets the status code for the add amazon RDS my SQL service default response
func (o *AddAmazonRDSMySQLServiceDefault) Code() int {
	return o._statusCode
}

func (o *AddAmazonRDSMySQLServiceDefault) Error() string {
	return fmt.Sprintf("[POST /v1/inventory/Services/AddAmazonRDSMySQL][%d] AddAmazonRDSMySQLService default  %+v", o._statusCode, o.Payload)
}

func (o *AddAmazonRDSMySQLServiceDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(AddAmazonRDSMySQLServiceDefaultBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*AddAmazonRDSMySQLServiceBody add amazon RDS my SQL service body
swagger:model AddAmazonRDSMySQLServiceBody
*/
type AddAmazonRDSMySQLServiceBody struct {

	// Instance endpoint (full DNS name). Required.
	Address string `json:"address,omitempty"`

	// Custom user-assigned labels.
	CustomLabels map[string]string `json:"custom_labels,omitempty"`

	// Environment name.
	Environment string `json:"environment,omitempty"`

	// Node identifier where this instance runs. Required.
	NodeID string `json:"node_id,omitempty"`

	// Instance port. Required.
	Port int64 `json:"port,omitempty"`

	// Unique across all Services user-defined name. Required.
	ServiceName string `json:"service_name,omitempty"`
}

// Validate validates this add amazon RDS my SQL service body
func (o *AddAmazonRDSMySQLServiceBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *AddAmazonRDSMySQLServiceBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *AddAmazonRDSMySQLServiceBody) UnmarshalBinary(b []byte) error {
	var res AddAmazonRDSMySQLServiceBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*AddAmazonRDSMySQLServiceDefaultBody ErrorResponse is a message returned on HTTP error.
swagger:model AddAmazonRDSMySQLServiceDefaultBody
*/
type AddAmazonRDSMySQLServiceDefaultBody struct {

	// code
	Code int32 `json:"code,omitempty"`

	// error
	Error string `json:"error,omitempty"`

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this add amazon RDS my SQL service default body
func (o *AddAmazonRDSMySQLServiceDefaultBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *AddAmazonRDSMySQLServiceDefaultBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *AddAmazonRDSMySQLServiceDefaultBody) UnmarshalBinary(b []byte) error {
	var res AddAmazonRDSMySQLServiceDefaultBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*AddAmazonRDSMySQLServiceOKBody add amazon RDS my SQL service OK body
swagger:model AddAmazonRDSMySQLServiceOKBody
*/
type AddAmazonRDSMySQLServiceOKBody struct {

	// amazon rds mysql
	AmazonRDSMysql *AddAmazonRDSMySQLServiceOKBodyAmazonRDSMysql `json:"amazon_rds_mysql,omitempty"`
}

// Validate validates this add amazon RDS my SQL service OK body
func (o *AddAmazonRDSMySQLServiceOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateAmazonRDSMysql(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *AddAmazonRDSMySQLServiceOKBody) validateAmazonRDSMysql(formats strfmt.Registry) error {

	if swag.IsZero(o.AmazonRDSMysql) { // not required
		return nil
	}

	if o.AmazonRDSMysql != nil {
		if err := o.AmazonRDSMysql.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("addAmazonRdsMySqlServiceOk" + "." + "amazon_rds_mysql")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *AddAmazonRDSMySQLServiceOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *AddAmazonRDSMySQLServiceOKBody) UnmarshalBinary(b []byte) error {
	var res AddAmazonRDSMySQLServiceOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*AddAmazonRDSMySQLServiceOKBodyAmazonRDSMysql AmazonRDSMySQLService represents a MySQL instance running on a single RemoteAmazonRDS Node.
swagger:model AddAmazonRDSMySQLServiceOKBodyAmazonRDSMysql
*/
type AddAmazonRDSMySQLServiceOKBodyAmazonRDSMysql struct {

	// Instance endpoint (full DNS name).
	Address string `json:"address,omitempty"`

	// Custom user-assigned labels.
	CustomLabels map[string]string `json:"custom_labels,omitempty"`

	// Environment name.
	Environment string `json:"environment,omitempty"`

	// Node identifier where this instance runs.
	NodeID string `json:"node_id,omitempty"`

	// Instance port.
	Port int64 `json:"port,omitempty"`

	// Unique randomly generated instance identifier.
	ServiceID string `json:"service_id,omitempty"`

	// Unique across all Services user-defined name.
	ServiceName string `json:"service_name,omitempty"`
}

// Validate validates this add amazon RDS my SQL service OK body amazon RDS mysql
func (o *AddAmazonRDSMySQLServiceOKBodyAmazonRDSMysql) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *AddAmazonRDSMySQLServiceOKBodyAmazonRDSMysql) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *AddAmazonRDSMySQLServiceOKBodyAmazonRDSMysql) UnmarshalBinary(b []byte) error {
	var res AddAmazonRDSMySQLServiceOKBodyAmazonRDSMysql
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
