package requestlog

import (
	"context"
	"errors"
	"time"

	"github.com/Hidayathamir/user-activity-tracking-go/internal/model"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/ctx/ctxclientauth"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/errkit"
	"github.com/Hidayathamir/user-activity-tracking-go/pkg/x"
	"github.com/redis/go-redis/v9"
)

// GetClientDailyRequestCount will get daily request count for last 7 day
// from redis. if not found will get from db and set cache to redis.
// if set cache to redis, will set ttl like this
//
//	for today,	 ttl 7 day
//	for today-1, ttl 6 day
//	for today-2, ttl 5 day
//	for today-3, ttl 4 day
//	for today-4, ttl 3 day
//	for today-5, ttl 2 day
//	for today-6, ttl 1 day
func (r *RequestLogUsecaseImpl) GetClientDailyRequestCount(ctx context.Context, req *model.ReqGetClientDailyRequestCount) (*model.ResGetClientDailyRequestCount, error) {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	clientAuth := ctxclientauth.Get(ctx)

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	dailyCountList := model.DailyCountList{}

	for i := range 7 {
		date := today.AddDate(0, 0, -i)
		ttl := time.Duration(7-i) * 24 * time.Hour

		count, err := r.getCountByAPIKeyAndDate(ctx, clientAuth.APIKey, date, ttl)
		if err != nil {
			return nil, errkit.AddFuncName(err)
		}

		dailyCountList = append(dailyCountList, model.DailyCount{
			Date:  date.Format(time.DateOnly),
			Count: count,
		})
	}

	res := &model.ResGetClientDailyRequestCount{
		DailyCountList: dailyCountList,
	}

	return res, nil
}

func (r *RequestLogUsecaseImpl) getCountByAPIKeyAndDate(ctx context.Context, apiKey string, datetime time.Time, ttl time.Duration) (int, error) {
	count, err := r.cache.GetCountByAPIKeyAndDate(ctx, apiKey, datetime)
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, errkit.AddFuncName(err)
	}

	if errors.Is(err, redis.Nil) {
		count, err = r.clientRequestCountRepository.GetCountByAPIKeyAndDate(ctx, r.db, apiKey, datetime)
		if err != nil {
			return 0, errkit.AddFuncName(err)
		}

		err = r.cache.SetClientRequestCount(ctx, apiKey, datetime, count, ttl)
		x.LogIfErrContext(ctx, err)
	}

	return count, nil
}
