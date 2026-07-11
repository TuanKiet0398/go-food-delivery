package restaurantstorage

import (
	"context"
	"food-delivery/common"
	restaurantmodel "food-delivery/module/restaurant/model"

	"gorm.io/gorm"
)



func (s *sqlStore) ListDataWithCondition(
	context context.Context,
	filter *restaurantmodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]restaurantmodel.Restaurant, error) {
	var result []restaurantmodel.Restaurant

	query := func() *gorm.DB {
		db := s.db.Where("status = ?", 1)

		if f := filter; f != nil {
			if f.OwnerId > 0 {
				db = db.Where("owner_id = ?", f.OwnerId)
			}
		}

		return db
	}

	if err := query().Table(restaurantmodel.Restaurant{}.TableName()).Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	paging.Fulfill()

	offset := (paging.Page - 1) * paging.Limit

	if err := query().
		Offset(offset).
		Limit(paging.Limit).
		Order("id desc").
		Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}