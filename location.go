package store

import (
	"context"
	"github.com/google/uuid"
	"github.com/r2day/db"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

const (
	storePos       = "store_position"
	defaultGEOUnit = "km"
	defaultRadius  = 10.0
)

var (
	geoUnit   = defaultGEOUnit
	geoRadius = defaultRadius
)

func init() {
	unit := os.Getenv("GEO_UNIT")
	if unit != "" {
		geoUnit = unit
	}

	radius := os.Getenv("GEO_RADIUS")
	if radius != "" {
		geoRadius, _ = strconv.ParseFloat(radius, 64)
	}
}

// GetStoresByLocation 获取门店id
// GEOSEARCH store_pos FROMLONLAT 13 38 BYRADIUS 100 km ASC
// GEOSEARCH store_position FROMLONLAT 114.03 23.54 BYRADIUS 7 km ASC
func GetStoresByLocation(ctx context.Context, longitude, latitude float64) []string {
	pos, err := db.RDB.GeoSearch(ctx, storePos, &redis.GeoSearchQuery{
		Longitude:  longitude,
		Latitude:   latitude,
		Sort:       "ASC", // 排序后返回最近的一个门店
		RadiusUnit: geoUnit,
		Radius:     geoRadius,
	}).Result()
	if err != nil {
		log.WithField("longitude", longitude).WithField("latitude", latitude).Error(err)
		return nil
	}
	return pos
}

// AddStore geoadd store_pos 13.361389 38.115556 "store_01" 15.087269 37.502669 "store_02"
func AddStore(ctx context.Context, stores []*Model) {
	for _, i := range stores {
		// 将key 记录下来以便退出的时候进行删除
		err := db.RDB.GeoAdd(ctx, storePos, &redis.GeoLocation{
			Name:      i.ID.Hex(),
			Longitude: i.Lbs.Longitude,
			Latitude:  i.Lbs.Latitude,
		}).Err()
		if err != nil {
			log.WithField("name", i.Name).Error(err)
			continue
		}
	}
}

// AddUserLocationForTemporary geoadd store_pos 13.361389 38.115556 "store_01" 15.087269 37.502669 "store_02"
func AddUserLocationForTemporary(ctx context.Context, long, lat float64) string {
	id := uuid.New()
	randUserLocationId := id.String()
	// 将key 记录下来以便退出的时候进行删除
	err := db.RDB.GeoAdd(ctx, storePos, &redis.GeoLocation{
		Name:      randUserLocationId,
		Longitude: long,
		Latitude:  lat,
	}).Err()
	if err != nil {
		log.WithField("name", randUserLocationId).Error(err)
	}
	return randUserLocationId
}

// GetDistance Sicily Palermo Catania km
// GetDistance geoadd store_pos 13.361389 38.115556 "store_01" 15.087269 37.502669 "store_02"
func GetDistance(ctx context.Context, member1, member2 string, uint string) float64 {
	dis, err := db.RDB.GeoDist(ctx, storePos, member1, member2, uint).Result()
	if err != nil {
		log.WithField("member1", member1).
			WithField("member2", member2).Error(err)
	}
	return dis
}

// RemoveUserTemporaryLocation Sicily Palermo Catania km
func RemoveUserTemporaryLocation(ctx context.Context, member string) int64 {
	counter, err := db.RDB.ZRem(ctx, storePos, member).Result()
	if err != nil {
		log.WithField("member", member).Error(err)
	}
	return counter
}
