package store

import (
	"github.com/open4go/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// CollectionNamePrefix 数据库表前缀
	// 可以根据具体业务的需要进行定义
	// 例如: sys_, scm_, customer_, order_ 等
	collectionNamePrefix = "mini_"
	// CollectionNameSuffix 后缀
	// 例如, _log, _config, _flow,
	collectionNameSuffix = "_manage"
	// 这个需要用户根据具体业务完成设定
	modelName = "store"
)

// LbsInfo 地址
type LbsInfo struct {
	Address   string  `json:"address" bson:"address"`
	Longitude float64 `json:"longitude"  bson:"longitude"`
	Latitude  float64 `json:"latitude"  bson:"latitude"`
	AreaName  string  `json:"areaName"  bson:"areaName"`
}

// Model 门店信息
type Model struct {
	// 模型继承
	model.Model `json:"-" bson:"-"`
	// 基本的数据库模型字段，一般情况所有model都应该包含如下字段
	// 创建时（用户上传的数据为空，所以默认可以不传该值)
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	// 分类/ 亦或则是分组等
	Category string `json:"category" bson:"category"`
	// BestMatch  最匹配的门店
	BestMatch bool `json:"best_match"`
	// 门店状态
	Status int `json:"status"  bson:"status"`
	// StoreID 门店id
	StoreID string `json:"storeId" bson:"storeId"`
	// Name 门店名称
	Name string `json:"name" bson:"name"`
	// Distance 根据用户当前位置计算出的结果
	Distance float64 `json:"distance" bson:"-"`
	// 门店营业时间
	ShopTime string `json:"shopTime" bson:"shopTime"`
	// 门店电话
	CallNumber string `json:"callNumber" bson:"callNumber"`
	// 地址
	Lbs LbsInfo `json:"lbs" bson:"lbs"`
}

// ResourceName 返回资源名称
func (m *Model) ResourceName() string {
	return modelName
}

// CollectionName 返回表名称
func (m *Model) CollectionName() string {
	return collectionNamePrefix + modelName + collectionNameSuffix
}
