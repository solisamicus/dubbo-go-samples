package service

import (
	"context"
	"encoding/json"

	"github.com/apache/dubbo-go-samples/transcation/seata-go/triple/proto"
	"github.com/dubbogo/gost/log/logger"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/seata/seata-go/pkg/rm/tcc"
	"github.com/seata/seata-go/pkg/tm"
	"google.golang.org/protobuf/types/known/anypb"
)

type UserProvider struct {
}

func (t *UserProvider) Prepare(ctx context.Context, params ...interface{}) (bool, error) {
	logger.Infof("Prepare result: %v, xid %v", params, tm.GetXID(ctx))
	return true, nil
}

func (t *UserProvider) Commit(ctx context.Context, businessActionContext *tm.BusinessActionContext) (bool, error) {
	logger.Infof("Commit result: %v, xid %s", businessActionContext, tm.GetXID(ctx))
	return true, nil
}

func (t *UserProvider) Rollback(ctx context.Context, businessActionContext *tm.BusinessActionContext) (bool, error) {
	logger.Infof("Rollback result: %v, xid %s", businessActionContext, tm.GetXID(ctx))
	return true, nil
}

func (t *UserProvider) GetActionName() string {
	logger.Infof("GetActionName result")
	return "TwoPhaseDemoService"
}

type UserProviderServer struct {
	*tcc.TCCServiceProxy
}

func (s *UserProviderServer) PrepareProxy(ctx context.Context, req *proto.PrepareRequest) (*proto.PrepareResponse, error) {
	logger.Info(tm.GetXID(ctx))
	ok, err := s.Prepare(ctx)
	return &proto.PrepareResponse{Result: ok.(bool)}, err
}

func (s *UserProviderServer) CommitProxy(ctx context.Context, req *proto.CommitRequest) (*proto.CommitResponse, error) {
	ok, err := s.Commit(ctx, &tm.BusinessActionContext{
		Xid:           req.BusinessActionContext.GetXid(),
		BranchId:      req.BusinessActionContext.GetBranchId(),
		ActionName:    req.BusinessActionContext.GetActionName(),
		ActionContext: convert(req.BusinessActionContext.GetActionContext()),
	})
	return &proto.CommitResponse{Result: ok}, err
}

func (s *UserProviderServer) RollbackProxy(ctx context.Context, req *proto.RollbackRequest) (*proto.RollbackResponse, error) {
	ok, err := s.Rollback(ctx, &tm.BusinessActionContext{
		Xid:           req.BusinessActionContext.GetXid(),
		BranchId:      req.BusinessActionContext.GetBranchId(),
		ActionName:    req.BusinessActionContext.GetActionName(),
		ActionContext: convert(req.BusinessActionContext.GetActionContext()),
	})
	return &proto.RollbackResponse{Result: ok}, err
}

func (s *UserProviderServer) GetActionNameProxy(ctx context.Context, req *proto.GetActionNameRequest) (*proto.GetActionNameResponse, error) {
	actionName := s.GetActionName()
	return &proto.GetActionNameResponse{ActionName: actionName}, nil
}

func convert(req map[string]*anypb.Any) map[string]interface{} {
	convertedMap := make(map[string]interface{})
	for key, value := range req {
		jsonStr, err := (&jsonpb.Marshaler{}).MarshalToString(value)
		if err != nil {
			logger.Error("converting marshal wrong")
		}

		var data interface{}
		err = json.Unmarshal([]byte(jsonStr), &data)
		if err != nil {
			logger.Error("converting unmarshal wrong")
		}
		convertedMap[key] = data
	}
	return convertedMap
}
