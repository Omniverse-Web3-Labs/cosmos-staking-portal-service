package model

const ValidatorAddressesTableName = "validator_addresses"

type ValidatorAddresses struct {
	ID          string `gorm:"column:id;type:varchar(255);not null;unique"`
	ValoperAddr string `gorm:"column:valoper_addr;type:varchar(255);not null;unique"`
	ValconsAddr string `gorm:"column:valcons_addr;type:varchar(255);not null;unique"`
}

func (v *ValidatorAddresses) TableName() string {
	return ValidatorAddressesTableName
}
