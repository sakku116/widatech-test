package repository

import (
	"backend/domain/dto"
	"backend/domain/model"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type InvoiceRepo struct {
	DB *gorm.DB
}

type IInvoiceRepo interface {
	Create(invoice *model.Invoice) error
	Update(invoice *model.Invoice) error
	Delete(invoice *model.Invoice) error
	GetByUUID(uuid string, preload bool) (*model.Invoice, error)
	GetList(
		params dto.InvoiceRepo_GetListParams,
	) ([]model.Invoice, int64, error)
}

func NewInvoiceRepo(db *gorm.DB) IInvoiceRepo {
	return &InvoiceRepo{
		DB: db,
	}
}

func (r *InvoiceRepo) Create(invoice *model.Invoice) error {
	err := r.DB.Create(invoice).Error
	if err != nil {
		return errors.New("failed to create invoice")
	}

	return err
}

func (r *InvoiceRepo) Update(invoice *model.Invoice) error {
	err := r.DB.Save(invoice).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("not found")
		}
		return errors.New("failed to get")
	}

	return err
}

func (r *InvoiceRepo) Delete(invoice *model.Invoice) error {
	err := r.DB.Delete(invoice).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("not found")
		}
		return errors.New("failed to get")
	}

	return err
}

func (r *InvoiceRepo) GetByUUID(uuid string, preload bool) (*model.Invoice, error) {
	var invoice model.Invoice
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

func (repo *InvoiceRepo) GetList(
	params dto.InvoiceRepo_GetListParams,
) ([]model.Invoice, int64, error) {
	var result []model.Invoice
	var totalData int64

	// validate param
	err := params.Validate()
	if err != nil {
		return result, totalData, err
	}

	tx := repo.DB.Model(&result)

	// filtering
	if params.PaymentType != nil {
		tx = tx.Where("payment_type = ?", *params.PaymentType)
	}
	if params.Query != nil {
		if params.QueryBy != nil {
			tx = tx.Where(fmt.Sprintf("%s LIKE ?", *params.QueryBy), "%"+*params.Query+"%")
		} else {
			// filter by all queriable fields
			conditions := ""
			conditionValues := []interface{}{}
			tmp := model.Invoice{}
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
