// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

// Package consul provides facilities for working with Consul.
package consul

import (
	"path"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

const (
	// prefix for all keys in KV operations.
	prefix = "percona/"
)

// FIXME remove alias after reconsidering the usage of GetNodes() method
// Currently it breaks encapsulation of this service
type Node = api.Node

// Client represents a client for Consul API.
// All keys in KV operations are prefixed to avoid collisions.
type Client struct {
	c *api.Client
}

// NewClient creates a new client for given Consul address.
func NewClient(addr string) (*Client, error) {
	c, err := api.NewClient(&api.Config{
		Address: addr,
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to Consul")
	}
	return &Client{c}, nil
}

// GetNodes returns list of nodes.
func (c *Client) GetNodes() ([]*Node, error) {
	nodes, _, err := c.c.Catalog().Nodes(nil)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get Nodes")
	}
	return nodes, nil
}

// GetNode returns a node
func (c *Client) GetNode(nodeID string) (*api.CatalogNode, error) {
	node, _, err := c.c.Catalog().Node(nodeID, nil)

	return node, err
}

// DeregisterNode deregisters consul catalog node
func (c *Client) DeregisterNode(nodeID string) (*api.WriteMeta, error) {
	return c.c.Catalog().Deregister(&api.CatalogDeregistration{
		Node: nodeID,
	}, nil)
}

// DeregisterService deregisters consul catalog service
func (c *Client) DeregisterService(nodeID, serviceID string) (*api.WriteMeta, error) {
	return c.c.Catalog().Deregister(&api.CatalogDeregistration{
		Node:      nodeID,
		ServiceID: serviceID,
	}, nil)
}

// GetKV returns value for a given key from Consul, or nil, if key does not exist.
func (c *Client) GetKV(key string) ([]byte, error) {
	key = path.Join(prefix, key)
	return c.getKV(key)
}

// GetExactKV returns value for a given key (without prefix) from Consul,
// or nil, if key does not exist.
func (c *Client) GetExactKV(key string) ([]byte, error) {
	return c.getKV(key)
}

func (c *Client) getKV(key string) ([]byte, error) {
	pair, _, err := c.c.KV().Get(key, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get key %q", key)
	}
	if pair == nil {
		return nil, nil
	}
	return pair.Value, nil
}

// PutKV puts given key/value pair into Consul.
func (c *Client) PutKV(key string, value []byte) error {
	key = path.Join(prefix, key)
	return c.putKV(key, value)
}

// PutExactKV puts given key/value pair into Consul.
func (c *Client) PutExactKV(key string, value []byte) error {
	return c.putKV(key, value)
}

func (c *Client) putKV(key string, value []byte) error {
	pair := &api.KVPair{Key: key, Value: value}
	_, err := c.c.KV().Put(pair, nil)
	return errors.Wrapf(err, "cannot put key %q", key)
}

// DeleteKV deletes given key from Consul.
func (c *Client) DeleteKV(key string) error {
	key = path.Join(prefix, key)
	_, err := c.c.KV().Delete(key, nil)
	return errors.Wrapf(err, "cannot delete key %q", key)
}

func (c *Client) DeleteKVTree(key string) error {
	_, err := c.c.KV().DeleteTree(key, nil)
	return errors.Wrapf(err, "cannot delete kv tree %q", key)
}
