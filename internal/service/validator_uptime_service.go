package service

import (
	"app/internal/db/postgres"
	dbModel "app/internal/db/postgres/model"
	"app/internal/model"
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cosmos/btcutil/bech32"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type validatorUptimeService struct {
}

var ValidatorUptimeService validatorUptimeService

// 获取验证器最近recentBlockCount个区块的在线情况
func (s *validatorUptimeService) ListRecentUptime(ctx context.Context, validatorAddressesHex []string, recentBlockCount int) ([]*model.ValidatorUptime, error) {
	subqueryDb := postgres.Subquery()
	//查询最后一个区块
	var lastBlockHeight dbModel.BlockEvents
	err := subqueryDb.Model(&dbModel.BlockEvents{}).WithContext(ctx).
		Select("block_height, timestamp").
		Order("timestamp DESC").
		Limit(1).Find(&lastBlockHeight).Error
	if err != nil {
		return nil, err
	}
	if lastBlockHeight.BlockHeight == 0 {
		return make([]*model.ValidatorUptime, 0), nil
	}
	//因为BlockSignatureRecords记录的是上一个区块的签名情况，所以需要从offset=limit个区块开始查询
	var startBlock dbModel.BlockEvents
	err = subqueryDb.Model(&dbModel.BlockEvents{}).WithContext(ctx).
		Select("block_height, timestamp").
		Where("timestamp <= ?", lastBlockHeight.Timestamp).
		Order("timestamp DESC").Offset(recentBlockCount).Limit(1).Find(&startBlock).Error
	if err != nil {
		return nil, err
	}
	if startBlock.BlockHeight == 0 {
		return make([]*model.ValidatorUptime, 0), nil
	}

	startTime := startBlock.Timestamp
	endTime := lastBlockHeight.Timestamp
	//查询验证器在limit个区块中的在线情况
	var validatorUptime []*model.ValidatorUptime
	err = subqueryDb.Model(&dbModel.BlockSignatureRecords{}).WithContext(ctx).
		Where("timestamp >= ? AND timestamp < ?", startTime, endTime).
		Where("validator IN (?)", validatorAddressesHex).
		Where("signed = true").
		Group("validator").Select("validator, count(*) as count").
		Find(&validatorUptime).Error
	if err != nil {
		return nil, err
	}

	return validatorUptime, nil
}

// 获取验证器最近recentBlockCount个区块的在线情况,validatorOperAddresses传入的是validatorCosmosvaloper地址
func (s *validatorUptimeService) ListRecentUptimeByOperAddresses(ctx context.Context, validatorOperAddresses []string, recentBlockCount int) ([]*model.ValidatorUptime, error) {
	if len(validatorOperAddresses) == 0 {
		return make([]*model.ValidatorUptime, 0), nil
	}
	validatorAddressesMapping, hexAddresses, err := s.BatchValOper2ValConsHex(ctx, validatorOperAddresses)

	if err != nil {
		return nil, err
	}

	list, err := s.ListRecentUptime(ctx, hexAddresses, recentBlockCount)
	if err != nil {
		return nil, err
	}
	for _, item := range list {
		item.Validator = validatorAddressesMapping[item.Validator]
	}
	return list, nil
}

// 批量将valoper地址转换为cons hex
func (s *validatorUptimeService) BatchValOper2ValConsHex(ctx context.Context, validatorOperAddresses []string) (map[string]string, []string, error) {

	subqueryDb := postgres.Subquery()

	validatorAddressesList := make([]*dbModel.ValidatorAddresses, 0)
	err := subqueryDb.Model(&dbModel.ValidatorAddresses{}).WithContext(ctx).
		Where("valoper_addr IN (?)", validatorOperAddresses).
		Find(&validatorAddressesList).Error
	if err != nil {
		return nil, nil, err
	}
	validatorAddressesHex := make([]string, 0)
	validatorAddressesHexMap := make(map[string]string)
	for _, validatorAddress := range validatorAddressesList {
		validatorAddressesHex = append(validatorAddressesHex, validatorAddress.ValconsAddr)
		validatorAddressesHexMap[validatorAddress.ValconsAddr] = validatorAddress.ValoperAddr
	}

	return validatorAddressesHexMap, validatorAddressesHex, nil
}

// valconsAddr to hex
func (s *validatorUptimeService) Bech32ToHex(bech32Addr string, prefix string) (string, error) {
	// 解码 Bech32 地址成 []byte
	addrBytes, err := sdk.GetFromBech32(bech32Addr, prefix)
	if err != nil {
		return "", err
	}
	// 转换为 HEX 字符串（大写）
	hexAddr := hex.EncodeToString(addrBytes)
	return strings.ToUpper(hexAddr), nil
}

func (s *validatorUptimeService) ConvertHexToBech32Address(hexAddr string, prefix string) (string, error) {
	// 去掉可能的 "0x" 前缀
	if len(hexAddr) >= 2 && hexAddr[:2] == "0x" {
		hexAddr = hexAddr[2:]
	}

	// 转换 hex 字符串为 byte 数组
	data, err := hex.DecodeString(hexAddr)
	if err != nil {
		return "", fmt.Errorf("invalid hex: %w", err)
	}

	// 将 8-bit bytes 转为 5-bit array（Bech32 需要）
	converted, err := bech32.ConvertBits(data, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("convert to 5-bit failed: %w", err)
	}

	// Bech32 编码
	bech32Addr, err := bech32.Encode(prefix, converted)
	if err != nil {
		return "", fmt.Errorf("bech32 encode failed: %w", err)
	}

	return bech32Addr, nil
}
