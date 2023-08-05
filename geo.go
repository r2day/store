package store

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

// FetchDataAndSaveOnCache 获取数据库并且存储到缓存中
func FetchDataAndSaveOnCache(c *gin.Context, longitude, latitude float64) (int64, *Model, []*Model, error) {
	// 标记最匹配的门店
	bestMatchStore := &Model{}
	allStores := make([]*Model, 0)
	// TODO 后期建设数据库的访问
	handler := bestMatchStore.Init(c.Request.Context(), MongoDatabase, bestMatchStore.CollectionName())
	counter, err := handler.GetList(bson.D{}, &allStores)
	if err != nil {
		return 0, bestMatchStore, allStores, err
	}

	// 存入redis 数据库
	AddStore(c.Request.Context(), allStores)
	// 获取最近的门店(列表)
	storeIDList := GetStoresByLocation(c.Request.Context(), longitude, latitude)

	// 同时存入用户位置到数据redis（便于下一步计算其到各个门店的距离）
	userLocationId := AddUserLocationForTemporary(c.Request.Context(), longitude, latitude)
	defer func() {
		RemoveUserTemporaryLocation(c.Request.Context(), userLocationId)
	}()

	// 取出最近的一个门店（用于直接返回给用户）
	storeId := ""
	if len(storeIDList) > 0 {
		storeId = storeIDList[0]
	}
	log.WithField("allStores", allStores).Warning("==========")
	storeCache := NewCache(c.Request.Context())
	for _, store := range allStores {
		if store.ID.Hex() == storeId {
			store.BestMatch = true
			bestMatchStore = store
		}
		// 计算距离(所有门店与客户的距离)
		dis := GetDistance(c.Request.Context(), store.ID.Hex(), userLocationId, "km")
		store.Distance = dis

		// 存入门店信息
		err := storeCache.Dump(store)
		if err != nil {
			log.WithField("store", store).Error(err)
			continue
		}
	}
	return counter, bestMatchStore, allStores, nil
}
