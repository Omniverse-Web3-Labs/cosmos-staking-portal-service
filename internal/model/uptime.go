package model

type ValidatorUptime struct {
	Validator string `json:"validator"`
	Count     int    `json:"count"`
}
