package middlewares

import (
	"context"

	"gitdev.devops.krungthai.com/aster/ariskill/app/user"
	"go.mongodb.org/mongo-driver/mongo"
)

// Both auth employee service and storage
type AuthEmployeeService interface {
	FindOne(ctx context.Context, filter any) (*user.Employee, error)
}

type authEmployeeService struct {
	coll *mongo.Collection
}

func NewAuthEmployeeService(coll *mongo.Collection) AuthEmployeeService {
	return &authEmployeeService{
		coll: coll,
	}
}

func (ae *authEmployeeService) FindOne(ctx context.Context, filter any) (*user.Employee, error) {
	res := &user.Employee{}
	err := ae.coll.FindOne(ctx, filter).Decode(res)
	return res, err
}
