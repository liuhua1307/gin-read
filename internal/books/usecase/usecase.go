package usecase

import (
	"github.com/google/wire"
	"github.com/liuhua1307/gin-read/internal/domain"
)

var ProviderSet = wire.NewSet(NewBookUseCase, wire.Bind(new(domain.BookUseCase), new(*BookUseCase)))
