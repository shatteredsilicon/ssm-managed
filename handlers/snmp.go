package handlers

import (
	"context"

	"github.com/shatteredsilicon/ssm-managed/api"
	"github.com/shatteredsilicon/ssm-managed/services/snmp"
	"github.com/shatteredsilicon/ssm-managed/utils/logger"
)

// SNMPServer handles requests to manage SNMP nodes and services.
type SNMPServer struct {
	SNMP *snmp.Service
}

// Add adds new SNMP instance.
func (s *SNMPServer) Add(ctx context.Context, req *api.SNMPAddRequest) (*api.SNMPAddResponse, error) {
	id, err := s.SNMP.Add(
		ctx, req.Name, req.Address, req.Port,
		req.Username, req.Password, req.Version,
		req.Community, req.SecurityLevel, req.AuthProtocol,
		req.PrivProtocol, req.PrivPassword, req.ContextName,
	)
	if err != nil {
		logger.Get(ctx).Errorf("%+v", err)
		return nil, err
	}

	resp := api.SNMPAddResponse{
		Id: id,
	}
	return &resp, nil
}

// List returns all SNMP instances
func (s *SNMPServer) List(ctx context.Context, req *api.SNMPListRequest) (*api.SNMPListResponse, error) {
	res, err := s.SNMP.List(ctx)
	if err != nil {
		logger.Get(ctx).Errorf("%+v", err)
		return nil, err
	}

	var resp api.SNMPListResponse
	for _, db := range res {
		resp.Instances = append(resp.Instances, &api.SNMPInstance{
			Node: &api.SNMPNode{
				Name: db.Node.Name,
			},
			Service: &api.SNMPService{
				Address:       *db.Service.Address,
				Port:          uint32(*db.Service.Port),
				Engine:        *db.Service.Engine,
				EngineVersion: *db.Service.EngineVersion,
			},
		})
	}
	return &resp, nil
}

// check interfaces
var (
	_ api.SNMPServer = (*SNMPServer)(nil)
)
