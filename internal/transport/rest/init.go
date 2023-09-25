package rest

import (
	"net/http"

	"github.com/Lalipopp4/test_api/internal/config"
	"github.com/Lalipopp4/test_api/internal/usecase/segment"
	"github.com/Lalipopp4/test_api/internal/usecase/user"
	"github.com/Lalipopp4/test_api/pkg/logging"
)

type restAPI struct {
	config         *config.Config
	logger         logging.Logger
	segmentService segment.SegmentService
	userService    user.UserService
	httpServer     *http.Server
}

func New() (Server, error) {
	cfg, err := config.InitConfig()
	if err != nil {
		return nil, err
	}
	logger, err := logging.New()
	if err != nil {
		return nil, err
	}
	segmentService, err := segment.New(cfg)
	if err != nil {
		return nil, err
	}
	userService, err := user.New(cfg)
	if err != nil {
		return nil, err
	}
	rAPI := &restAPI{
		segmentService: segmentService,
		userService:    userService,
		logger:         logger,
		config:         cfg,
	}
	rAPI.httpServer = &http.Server{
		Addr:    cfg.Server.Addr + cfg.Server.Port,
		Handler: handle(rAPI),
	}
	return rAPI, nil
}
