package handlers

import (
	"context"

	"github.com/shatteredsilicon/ssm-managed/api"
	"github.com/shatteredsilicon/ssm-managed/services/mongodb"
	"github.com/shatteredsilicon/ssm-managed/utils/logger"
)

// MongoDBServer handles requests to manage MongoDB nodes and services.
type MongoDBServer struct {
	MongoDB *mongodb.Service
}

// Add adds new MongoDB instance.
func (s *MongoDBServer) Add(ctx context.Context, req *api.MongoDBAddRequest) (*api.MongoDBAddResponse, error) {
	id, err := s.MongoDB.Add(ctx, req.Name, req.Address, nil)
	if err != nil {
		logger.Get(ctx).Errorf("%+v", err)
		return nil, err
	}

	resp := api.MongoDBAddResponse{
		Id: id,
	}
	return &resp, nil
}

// check interfaces
var (
	_ api.MongoDBServer = (*MongoDBServer)(nil)
)
