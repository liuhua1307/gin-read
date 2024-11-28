package repository

import (
	"github.com/google/wire"
	"github.com/liuhua1307/gin-read/internal/domain"
)

var ProviderSet = wire.NewSet(NewBookMySQLRepository, wire.Bind(new(domain.BookRepository), new(*BookMySQLRepository)))
