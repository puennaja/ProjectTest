package auth

import (
	"io/ioutil"
	"os"
	"ticket/internal/core/domain"
	"ticket/internal/core/port"

	"github.com/casbin/casbin/v2"
	jsonadapter "github.com/casbin/json-adapter/v2"
	"go.uber.org/zap"
)

type Config struct {
	Config  domain.AuthConfig
	UserSrv port.UserService
}

type Service struct {
	logger         *zap.SugaredLogger
	config         domain.AuthConfig
	userSrv        port.UserService
	casbinEnforcer casbin.IEnforcer
}

func New(logger *zap.SugaredLogger, cfg Config) port.AuthService {
	jsonFile, err := os.Open(cfg.Config.PolicyPath)
	if err != nil {
		logger.Panic("Auth Service [New]: " + err.Error())
	}
	defer jsonFile.Close()

	jsonByte, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logger.Panic("Auth Service [New]: " + err.Error())
	}

	enforcer, err := casbin.NewEnforcer(cfg.Config.ModelPath, jsonadapter.NewAdapter(&jsonByte))
	if err != nil {
		logger.Panic("Auth Service [New]: " + err.Error())
	}
	err = enforcer.LoadPolicy()
	if err != nil {
		logger.Error("Auth Service [New]: " + err.Error())
	}
	return &Service{
		logger:         logger,
		config:         cfg.Config,
		userSrv:        cfg.UserSrv,
		casbinEnforcer: enforcer,
	}
}
