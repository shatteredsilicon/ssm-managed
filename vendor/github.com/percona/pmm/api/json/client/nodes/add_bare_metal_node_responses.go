// Code generated by go-swagger; DO NOT EDIT.

package nodes

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/percona/pmm/api/json/models"
)

// AddBareMetalNodeReader is a Reader for the AddBareMetalNode structure.
type AddBareMetalNodeReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AddBareMetalNodeReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewAddBareMetalNodeOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewAddBareMetalNodeOK creates a AddBareMetalNodeOK with default headers values
func NewAddBareMetalNodeOK() *AddBareMetalNodeOK {
	return &AddBareMetalNodeOK{}
}

/*AddBareMetalNodeOK handles this case with default header values.

(empty)
*/
type AddBareMetalNodeOK struct {
	Payload *models.InventoryAddBareMetalNodeResponse
}

func (o *AddBareMetalNodeOK) Error() string {
	return fmt.Sprintf("[POST /v0/inventory/Nodes/AddBareMetal][%d] addBareMetalNodeOK  %+v", 200, o.Payload)
}

func (o *AddBareMetalNodeOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.InventoryAddBareMetalNodeResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
