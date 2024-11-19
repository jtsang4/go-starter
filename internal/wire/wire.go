package wire

import (
	"github.com/google/wire"
	"github.com/jtsang4/go-stater/config"
	"github.com/jtsang4/go-stater/internal/api"
	"github.com/jtsang4/go-stater/internal/repository"
	"github.com/jtsang4/go-stater/internal/service"
	"github.com/jtsang4/go-stater/pkg/cache"
	"github.com/jtsang4/go-stater/pkg/database"
	"gorm.io/gorm"
)

// ProviderSet 是所有provider的集合
var ProviderSet = wire.NewSet(
	ProvideUserRepository,
	ProvideUserService,
	ProvideUserHandler,
)

func ProvideUserRepository(db *gorm.DB) *repository.UserRepository {
	return repository.NewUserRepository(db)
}

func ProvideUserService(repo *repository.UserRepository, cache *cache.RedisCache) *service.UserService {
	return service.NewUserService(repo, cache)
}

func ProvideUserHandler(s *service.UserService, cfg *config.Config) *api.UserHandler {
	return api.NewUserHandler(s, cfg)
}

// InitializeAPI 初始化API服务
func InitializeAPI(cfg *config.Config) (*api.UserHandler, error) {
	wire.Build(
		database.InitDB,
		ProviderSet,
	)
	return nil, nil
}
