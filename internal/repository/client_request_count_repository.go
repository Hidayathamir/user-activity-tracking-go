package repository

import (
	"context"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/column"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/table"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	sq "github.com/Masterminds/squirrel"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/MockRepositoryClientRequestCount.go -pkg=mock . ClientRequestCountRepository

type ClientRequestCountRepository interface {
	IncrementCount(ctx context.Context, db *gorm.DB, apiKey string, datetime time.Time, count int) (int, error)
}

var _ ClientRequestCountRepository = &ClientRequestCountRepositoryImpl{}

type ClientRequestCountRepositoryImpl struct {
	Config *viper.Viper
}

func NewClientRequestCountRepository(cfg *viper.Viper) *ClientRequestCountRepositoryImpl {
	return &ClientRequestCountRepositoryImpl{
		Config: cfg,
	}
}

// IncrementCount increments count for a given apiKey and datetime.
//
// Behavior:
//   - Inserts a new row if it does not exist.
//   - If (api_key, datetime) already exists, increments count by argument `count`.
//   - Returns the latest updated count.
//
// SQL behavior:
//
//	INSERT ... ON CONFLICT ... DO UPDATE SET count = count + {count} RETURNING count
func (r *ClientRequestCountRepositoryImpl) IncrementCount(ctx context.Context, db *gorm.DB, apiKey string, datetime time.Time, count int) (int, error) {
	sqlStr, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(table.ClientRequestCount).
		Columns(column.APIKey.Str(), column.Datetime.Str(), column.Count.Str()).
		Values(apiKey, datetime, count).
		Suffix(`
            ON CONFLICT (`+column.APIKey.Str()+`, `+column.Datetime.Str()+`)
            DO UPDATE SET `+column.Count.Str()+` = `+table.ClientRequestCount+`.`+column.Count.Str()+` + ?
            RETURNING `+column.Count.Str()+`
        `, count).
		ToSql()
	if err != nil {
		return 0, errkit.AddFuncName(err)
	}

	var newCount int
	err = db.WithContext(ctx).Raw(sqlStr, args...).Scan(&newCount).Error
	if err != nil {
		return 0, errkit.AddFuncName(err)
	}

	return newCount, nil
}
