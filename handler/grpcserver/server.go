package grpcserver

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/n-creativesystem/rbns/infra"
	"github.com/n-creativesystem/rbns/infra/dao"
	"github.com/n-creativesystem/rbns/logger"
	"github.com/n-creativesystem/rbns/proto"
	"github.com/n-creativesystem/rbns/service"
	"google.golang.org/grpc"
)

type Option func(*config)

type config struct {
	secure bool
}

func WithSecure(conf *config) {
	conf.secure = true
}

func New(db dao.DataBase, opts ...Option) *grpc.Server {
	conf := &config{}
	for _, opt := range opts {
		opt(conf)
	}
	interceptors := []grpc.UnaryServerInterceptor{
		logger.GrpcLogger(),
		Recovery,
	}
	repo := infra.NewRepository(db)
	if conf.secure {
		authRepo := infra.NewAuth(db)
		token := authRepo.Get()
		if token == "" {
			token_, err := authRepo.Generate()
			if err != nil {
				panic(err)
			}
			token = token_
		}
		interceptors = append(interceptors, AuthUnaryServerInterceptor(token))
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(interceptors...),
		),
	)
	pSrv := service.NewPermissionService(repo)
	rSrv := service.NewRoleService(repo)
	oSrv := service.NewOrganizationService(repo)
	uSrv := service.NewUserService(repo)
	proto.RegisterPermissionServer(server, pSrv)
	proto.RegisterRoleServer(server, rSrv)
	proto.RegisterOrganizationServer(server, oSrv)
	proto.RegisterUserServer(server, uSrv)
	return server
}
