package account

import "gorm.io/gorm"

type AccountRepository struct {
	db *gorm.DB
}

type SAccount struct {
	Username string `json:"username",form:"username"`
	Password string `json:"password",form:"password"`
	Admin    int8   `json:"admin",form:"admin"`
	Tenant   string `json:"tenant",form:"tenant"`
	Token    string `json:"token",form:"token"`
}

type SPage struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type SPAccount struct {
	SAccount
	SPage
}

func (SAccount) TableName() string {
	return "account"
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (resp *AccountRepository) Query(query SPAccount) []SAccount {
	accounts := make([]SAccount, 0)
	resp.db.Find(&accounts, query)
	return accounts
}

func (resp *AccountRepository) Create(entity SAccount) (int64, error) {
	result := resp.db.Create(&entity)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (resp *AccountRepository) Delete(value SAccount) (int64, error) {
	var entity SAccount
	// resp.db.First(&entity, &value)
	result := resp.db.Where(&value).Delete(&entity)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (resp *AccountRepository) Update(entity SAccount, where SPAccount) (int64, error) {
	result := resp.db.Model(&SAccount{}).Updates(entity).Where(where)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
