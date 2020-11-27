// Code generated by go-swagger; DO NOT EDIT.

package channels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new channels API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for channels API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientService is the interface for Client methods
type ClientService interface {
	AddChannel(params *AddChannelParams) (*AddChannelOK, error)

	ChangeChannel(params *ChangeChannelParams) (*ChangeChannelOK, error)

	ListChannels(params *ListChannelsParams) (*ListChannelsOK, error)

	RemoveChannel(params *RemoveChannelParams) (*RemoveChannelOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  AddChannel adds channel adds notification channel
*/
func (a *Client) AddChannel(params *AddChannelParams) (*AddChannelOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewAddChannelParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "AddChannel",
		Method:             "POST",
		PathPattern:        "/v1/management/ia/Channels/Add",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &AddChannelReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*AddChannelOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*AddChannelDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ChangeChannel changes channel changes notification channel
*/
func (a *Client) ChangeChannel(params *ChangeChannelParams) (*ChangeChannelOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewChangeChannelParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "ChangeChannel",
		Method:             "POST",
		PathPattern:        "/v1/management/ia/Channels/Change",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ChangeChannelReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ChangeChannelOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ChangeChannelDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ListChannels lists channels returns a list of all notifation channels
*/
func (a *Client) ListChannels(params *ListChannelsParams) (*ListChannelsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListChannelsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "ListChannels",
		Method:             "POST",
		PathPattern:        "/v1/management/ia/Channels/List",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ListChannelsReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListChannelsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListChannelsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  RemoveChannel removes channel removes notification channel
*/
func (a *Client) RemoveChannel(params *RemoveChannelParams) (*RemoveChannelOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewRemoveChannelParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "RemoveChannel",
		Method:             "POST",
		PathPattern:        "/v1/management/ia/Channels/Remove",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &RemoveChannelReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*RemoveChannelOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*RemoveChannelDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
