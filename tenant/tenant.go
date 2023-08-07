package tenant

import "gorm.io/gorm"

type TenantRepository struct {
	db *gorm.DB
}

type STenant struct {
	Name        string `json:"name",	form:"name"`
	Owner       string `json:"owner",	form:"owner"`
	Description string `json:"description",form:"description"`
}

type SPage struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type SPTelant struct {
	STenant
	SPage
}

func (STenant) TableName() string {
	return "tenant"
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{
		db: db,
	}
}

func (resp *TenantRepository) Query(query SPTelant) []STenant {
	tenants := make([]STenant, 0)
	resp.db.Find(&tenants, query)
	return tenants
}

func (resp *TenantRepository) Create(entity STenant) (int64, error) {
	result := resp.db.Create(entity)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (resp *TenantRepository) Delete(value STenant) (int64, error) {
	var entity STenant
	// resp.db.First(&entity, &value)
	result := resp.db.Where(&value).Delete(&entity)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (resp *TenantRepository) Update(entity STenant, where SPTelant) (int64, error) {
	result := resp.db.Model(&STenant{}).Where(where).Updates(entity)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
