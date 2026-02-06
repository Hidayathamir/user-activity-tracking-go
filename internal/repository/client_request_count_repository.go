package repository

import (
	"context"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/converter"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/entity"
	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/column"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/constant/table"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/dbretry"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	sq "github.com/Masterminds/squirrel"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../mock/MockRepositoryClientRequestCount.go -pkg=mock . ClientRequestCountRepository

type ClientRequestCountRepository interface {
	IncrementCount(ctx context.Context, db *gorm.DB, apiKey string, datetime time.Time, count int) (int, error)
	GetTop3ClientRequestCount24Hour(ctx context.Context, db *gorm.DB) (model.APIKeyCountList, error)
	GetCountByAPIKeyAndDate(ctx context.Context, db *gorm.DB, apiKey string, datetime time.Time) (int, error)
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
	err = dbretry.Do(func() error {
		return db.WithContext(ctx).Raw(sqlStr, args...).Scan(&newCount).Error
	})
	if err != nil {
		return 0, errkit.AddFuncName(err)
	}

	return newCount, nil
}

func (r *ClientRequestCountRepositoryImpl) GetTop3ClientRequestCount24Hour(ctx context.Context, db *gorm.DB) (model.APIKeyCountList, error) {
	var top3ClientRequestCountList entity.Top3ClientRequestCountList

	_24HourAgo := time.Now().UTC().Add(-24 * time.Hour)

	err := dbretry.Do(func() error {
		return db.WithContext(ctx).
			Table(table.ClientRequestCount).
			Select([]string{column.APIKey.Str(), column.Count.SumAs("total_sum")}).
			Where(column.Datetime.GTE(_24HourAgo)).
			Group(column.APIKey.Str()).
			Order("total_sum DESC").
			Limit(3).
			Scan(&top3ClientRequestCountList).Error
	})
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := model.APIKeyCountList{}
	converter.EntityTop3ClientRequestCountListToModelAPIKeyCountList(&top3ClientRequestCountList, &res)

	return res, nil
}

func (r *ClientRequestCountRepositoryImpl) GetCountByAPIKeyAndDate(ctx context.Context, db *gorm.DB, apiKey string, datetime time.Time) (int, error) {
	startOfDay := time.Date(datetime.Year(), datetime.Month(), datetime.Day(), 0, 0, 0, 0, datetime.Location())
	startOfNextDay := startOfDay.AddDate(0, 0, 1)

	var count int
	err := dbretry.Do(func() error {
		return db.WithContext(ctx).
			Table(table.ClientRequestCount).
			Select("COALESCE(SUM(" + column.Count.Str() + "), 0)").
			Where(column.APIKey.Eq(apiKey)).
			Where(column.Datetime.GTE(startOfDay)).
			Where(column.Datetime.LT(startOfNextDay)).
			Scan(&count).Error
	})
	if err != nil {
		return 0, errkit.AddFuncName(err)
	}

	return count, nil
}
