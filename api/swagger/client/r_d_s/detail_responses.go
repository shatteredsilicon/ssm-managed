// Code generated by go-swagger; DO NOT EDIT.

package r_d_s

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/shatteredsilicon/ssm-managed/api/swagger/models"
)

// DetailReader is a Reader for the Detail structure.
type DetailReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DetailReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewDetailOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewDetailOK creates a DetailOK with default headers values
func NewDetailOK() *DetailOK {
	return &DetailOK{}
}

/*DetailOK handles this case with default header values.

(empty)
*/
type DetailOK struct {
	Payload *models.APIRDSDetailResponse
}

func (o *DetailOK) Error() string {
	return fmt.Sprintf("[GET /v0/rds/detail][%d] detailOK  %+v", 200, o.Payload)
}

func (o *DetailOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIRDSDetailResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
