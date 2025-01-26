package repository

import (
	"backend/domain/dto"
	"backend/domain/model"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type ProductRepo struct {
	DB *gorm.DB
}

type IProductRepo interface {
	Create(invoice *model.Product) error
	Update(invoice *model.Product) error
	Delete(invoice *model.Product) error
	GetByUUID(uuid string, preload bool) (*model.Product, error)
	GetList(
		params dto.ProductRepo_GetListParams,
	) ([]model.Product, int64, error)
	GetListByUUIDs(uuids []string) ([]model.Product, error)
}

func NewProductRepo(db *gorm.DB) IProductRepo {
	return &ProductRepo{
		DB: db,
	}
}

func (r *ProductRepo) Create(invoice *model.Product) error {
	err := r.DB.Create(invoice).Error
	if err != nil {
		return errors.New("failed to create invoice")
	}

	return err
}

func (r *ProductRepo) Update(invoice *model.Product) error {
	err := r.DB.Save(invoice).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("not found")
		}
		return errors.New("failed to get")
	}

	return err
}

func (r *ProductRepo) Delete(invoice *model.Product) error {
	err := r.DB.Delete(invoice).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("not found")
		}
		return errors.New("failed to get")
	}

	return err
}

func (r *ProductRepo) GetByUUID(uuid string, preload bool) (*model.Product, error) {
	var invoice model.Product
	query := r.DB
	if preload {
		query = query.Preload("Products")
	}
	err := query.Where("uuid = ?", uuid).First(&invoice).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("not found")
		}
		return nil, errors.New("failed to get")
	}

	return &invoice, err
}

func (repo *ProductRepo) GetListByUUIDs(uuids []string) ([]model.Product, error) {
	var result []model.Product
	query := repo.DB
	err := query.Where("uuid IN ?", uuids).Find(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("not found")
		}
		return nil, errors.New("failed to get")
	}

	return result, err
}

func (repo *ProductRepo) GetList(
	params dto.ProductRepo_GetListParams,
) ([]model.Product, int64, error) {
	var result []model.Product
	var totalData int64

	// validate param
	err := params.Validate()
	if err != nil {
		return result, totalData, err
	}

	tx := repo.DB.Model(&result)

	// filtering
	if params.Query != nil {
		if params.QueryBy != nil {
			tx = tx.Where(fmt.Sprintf("%s LIKE ?", *params.QueryBy), "%"+*params.Query+"%")
		} else {
			// filter by all queriable fields
			conditions := ""
			conditionValues := []interface{}{}
			tmp := model.Product{}
			queriableFields := tmp.GetProps().QueriableFields
			for _, field := range queriableFields {
				logger.Debugf("field: %s", field)
				if field == "" {
					logger.Debugf("skipping empty field")
					continue
				}
				conditions += fmt.Sprintf(
					`%s LIKE ? OR `,
					field,
				)
				conditionValues = append(conditionValues, "%"+*params.Query+"%")
			}
			logger.Debugf("conditionValues: %v", conditionValues)
			conditions = strings.TrimSuffix(conditions, " OR ")
			tx = tx.Where(
				conditions,
				conditionValues...,
			)
		}
	}

	// get count if needed
	if params.DoCount {
		err = tx.Count(&totalData).Error
		if err != nil {
			return nil, totalData, errors.New("failed to count: " + err.Error())
		}
	}

	// sorting
	if params.SortOrder != nil && params.SortBy != nil {
		tx = tx.Order(fmt.Sprintf("%s %s", *params.SortBy, *params.SortOrder))
	}

	// pagination
	if params.Page != nil && params.Limit != nil {
		tx = tx.Offset((*params.Page - 1) * *params.Limit).Limit(*params.Limit)
	}

	err = tx.Find(&result).Error
	if err != nil {
		return nil, totalData, errors.New("failed to get: " + err.Error())
	}

	return result, totalData, nil
}
