package protocol

import (
	"log"
	"ticket/config"
	"ticket/infrastructure"
	"ticket/internal/repository"
	"ticket/pkg/validator"

	"ticket/internal/core/domain"
	"ticket/internal/core/port"
	authSrv "ticket/internal/core/service/auth"
	ticketSrv "ticket/internal/core/service/ticket"
	userSrv "ticket/internal/core/service/user"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var app *application

type application struct {
	logger *zap.SugaredLogger
	svc    service
	pkg    packages
}

type service struct {
	authSvc   port.AuthService
	ticketSvc port.TicketService
	userSvc   port.UserService
}

type packages struct {
	validator validator.Validator
}

func init() {
	startApp()
}

func startApp() {
	// Setup
	config.Init()
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	loggerConfig.InitialFields = map[string]interface{}{
		"service-name":    config.GetConfig().AppConfig.Service,
		"service-version": config.GetConfig().AppConfig.Version,
	}
	defaultLog, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	logger := defaultLog.Sugar()

	// infrastructure
	var (
		mongoConn = infrastructure.InitMongo()
	)

	// packages
	var (
		trans     = validator.NewTranslator()
		validator = validator.New(trans)
		packages  = packages{
			validator: validator,
		}
	)

	// adapters
	var ()

	// repositories
	var (
		userRepo          = repository.NewUserRepository(logger, mongoConn, config.GetConfig().MongoDB.Database)
		ticketRepo        = repository.NewTicketRepository(logger, mongoConn, config.GetConfig().MongoDB.Database)
		ticketHistoryRepo = repository.NewTicketHistoryRepository(logger, mongoConn, config.GetConfig().MongoDB.Database)
		ticketCommentRepo = repository.NewTicketCommentRepository(logger, mongoConn, config.GetConfig().MongoDB.Database)
	)

	// service
	var (
		userSvc = userSrv.New(logger, userSrv.Config{
			UserRepo: userRepo,
		})
		authSvc = authSrv.New(logger, authSrv.Config{
			Config: domain.AuthConfig{
				AccessTokenExpire:  config.GetConfig().Service.AuthSvc.AccessTokenExpire,
				RefreshTokenExpire: config.GetConfig().Service.AuthSvc.RefreshTokenExpire,
				Secret:             config.GetConfig().Service.AuthSvc.Secret,
				ModelPath:          config.GetConfig().Service.AuthSvc.ModelPath,
				PolicyPath:         config.GetConfig().Service.AuthSvc.PolicyPath,
			},
			UserSrv: userSvc,
		})
		ticketSvc = ticketSrv.New(logger, ticketSrv.Config{
			TicketRepo:        ticketRepo,
			TicketHistoryRepo: ticketHistoryRepo,
			TicketCommentRepo: ticketCommentRepo,
		})
	)

	app = &application{
		logger: logger,
		svc: service{
			userSvc:   userSvc,
			authSvc:   authSvc,
			ticketSvc: ticketSvc,
		},
		pkg: packages,
	}
}
