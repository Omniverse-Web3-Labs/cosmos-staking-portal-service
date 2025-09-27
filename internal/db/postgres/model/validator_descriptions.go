package model

const ValidatorDescriptionsTableName = "validator_descriptions"

type ValidatorDescriptions struct {
	ID              string `gorm:"column:id;type:varchar(255);not null;unique"`
	Moniker         string `gorm:"column:moniker;type:varchar(255);not null"`
	Identity        string `gorm:"column:identity;type:varchar(255);not null"`
	Website         string `gorm:"column:website;type:varchar(255);not null"`
	SecurityContact string `gorm:"column:security_contact;type:varchar(255);not null"`
	Details         string `gorm:"column:details;type:varchar(255);not null"`
	Avatar          string `gorm:"column:avatar;type:varchar(255);not null"`
	ExtInfo         string `gorm:"column:ext_info;type:varchar(255);not null"`
	Timestamp       int64  `gorm:"column:timestamp;type:bigint;not null"`
}

func (v *ValidatorDescriptions) TableName() string {
	return ValidatorDescriptionsTableName
}
