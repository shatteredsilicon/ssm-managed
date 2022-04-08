// Code generated by go-swagger; DO NOT EDIT.

package scrape_configs

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

	models "github.com/shatteredsilicon/ssm-managed/api/swagger/models"
)

// NewCreateMixin7Params creates a new CreateMixin7Params object
// with the default values initialized.
func NewCreateMixin7Params() *CreateMixin7Params {
	var ()
	return &CreateMixin7Params{

		timeout: cr.DefaultTimeout,
	}
}

// NewCreateMixin7ParamsWithTimeout creates a new CreateMixin7Params object
// with the default values initialized, and the ability to set a timeout on a request
func NewCreateMixin7ParamsWithTimeout(timeout time.Duration) *CreateMixin7Params {
	var ()
	return &CreateMixin7Params{

		timeout: timeout,
	}
}

// NewCreateMixin7ParamsWithContext creates a new CreateMixin7Params object
// with the default values initialized, and the ability to set a context for a request
func NewCreateMixin7ParamsWithContext(ctx context.Context) *CreateMixin7Params {
	var ()
	return &CreateMixin7Params{

		Context: ctx,
	}
}

// NewCreateMixin7ParamsWithHTTPClient creates a new CreateMixin7Params object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewCreateMixin7ParamsWithHTTPClient(client *http.Client) *CreateMixin7Params {
	var ()
	return &CreateMixin7Params{
		HTTPClient: client,
	}
}

/*CreateMixin7Params contains all the parameters to send to the API endpoint
for the create mixin7 operation typically these are written to a http.Request
*/
type CreateMixin7Params struct {

	/*Body*/
	Body *models.APIScrapeConfigsCreateRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the create mixin7 params
func (o *CreateMixin7Params) WithTimeout(timeout time.Duration) *CreateMixin7Params {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create mixin7 params
func (o *CreateMixin7Params) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create mixin7 params
func (o *CreateMixin7Params) WithContext(ctx context.Context) *CreateMixin7Params {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create mixin7 params
func (o *CreateMixin7Params) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create mixin7 params
func (o *CreateMixin7Params) WithHTTPClient(client *http.Client) *CreateMixin7Params {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create mixin7 params
func (o *CreateMixin7Params) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the create mixin7 params
func (o *CreateMixin7Params) WithBody(body *models.APIScrapeConfigsCreateRequest) *CreateMixin7Params {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create mixin7 params
func (o *CreateMixin7Params) SetBody(body *models.APIScrapeConfigsCreateRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *CreateMixin7Params) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
