// Code generated by go-swagger; DO NOT EDIT.

package r_d_s

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDetailParams creates a new DetailParams object
// with the default values initialized.
func NewDetailParams() *DetailParams {
	var ()
	return &DetailParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDetailParamsWithTimeout creates a new DetailParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDetailParamsWithTimeout(timeout time.Duration) *DetailParams {
	var ()
	return &DetailParams{

		timeout: timeout,
	}
}

// NewDetailParamsWithContext creates a new DetailParams object
// with the default values initialized, and the ability to set a context for a request
func NewDetailParamsWithContext(ctx context.Context) *DetailParams {
	var ()
	return &DetailParams{

		Context: ctx,
	}
}

// NewDetailParamsWithHTTPClient creates a new DetailParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDetailParamsWithHTTPClient(client *http.Client) *DetailParams {
	var ()
	return &DetailParams{
		HTTPClient: client,
	}
}

/*DetailParams contains all the parameters to send to the API endpoint
for the detail operation typically these are written to a http.Request
*/
type DetailParams struct {

	/*QanDbInstanceUUID*/
	QanDbInstanceUUID *string
	/*ServerPassword*/
	ServerPassword *string
	/*ServerUsername*/
	ServerUsername *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the detail params
func (o *DetailParams) WithTimeout(timeout time.Duration) *DetailParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the detail params
func (o *DetailParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the detail params
func (o *DetailParams) WithContext(ctx context.Context) *DetailParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the detail params
func (o *DetailParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the detail params
func (o *DetailParams) WithHTTPClient(client *http.Client) *DetailParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the detail params
func (o *DetailParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithQanDbInstanceUUID adds the qanDbInstanceUUID to the detail params
func (o *DetailParams) WithQanDbInstanceUUID(qanDbInstanceUUID *string) *DetailParams {
	o.SetQanDbInstanceUUID(qanDbInstanceUUID)
	return o
}

// SetQanDbInstanceUUID adds the qanDbInstanceUuid to the detail params
func (o *DetailParams) SetQanDbInstanceUUID(qanDbInstanceUUID *string) {
	o.QanDbInstanceUUID = qanDbInstanceUUID
}

// WithServerPassword adds the serverPassword to the detail params
func (o *DetailParams) WithServerPassword(serverPassword *string) *DetailParams {
	o.SetServerPassword(serverPassword)
	return o
}

// SetServerPassword adds the serverPassword to the detail params
func (o *DetailParams) SetServerPassword(serverPassword *string) {
	o.ServerPassword = serverPassword
}

// WithServerUsername adds the serverUsername to the detail params
func (o *DetailParams) WithServerUsername(serverUsername *string) *DetailParams {
	o.SetServerUsername(serverUsername)
	return o
}

// SetServerUsername adds the serverUsername to the detail params
func (o *DetailParams) SetServerUsername(serverUsername *string) {
	o.ServerUsername = serverUsername
}

// WriteToRequest writes these params to a swagger request
func (o *DetailParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.QanDbInstanceUUID != nil {

		// query param qan_db_instance_uuid
		var qrQanDbInstanceUUID string
		if o.QanDbInstanceUUID != nil {
			qrQanDbInstanceUUID = *o.QanDbInstanceUUID
		}
		qQanDbInstanceUUID := qrQanDbInstanceUUID
		if qQanDbInstanceUUID != "" {
			if err := r.SetQueryParam("qan_db_instance_uuid", qQanDbInstanceUUID); err != nil {
				return err
			}
		}

	}

	if o.ServerPassword != nil {

		// query param server_password
		var qrServerPassword string
		if o.ServerPassword != nil {
			qrServerPassword = *o.ServerPassword
		}
		qServerPassword := qrServerPassword
		if qServerPassword != "" {
			if err := r.SetQueryParam("server_password", qServerPassword); err != nil {
				return err
			}
		}

	}

	if o.ServerUsername != nil {

		// query param server_username
		var qrServerUsername string
		if o.ServerUsername != nil {
			qrServerUsername = *o.ServerUsername
		}
		qServerUsername := qrServerUsername
		if qServerUsername != "" {
			if err := r.SetQueryParam("server_username", qServerUsername); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
