package service

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestValidatorUptimeService_ListValidatorUptime(t *testing.T) {

	validatorAddresses := []string{"7DE34C6F99726CC46C133C78DF6B5F714806EF4B"}
	validatorUptime, err := ValidatorUptimeService.ListRecentUptime(context.Background(), validatorAddresses, 100)
	if err != nil {
		t.Fatalf("failed to get validator uptime: %v", err)
	}
	fmt.Println(validatorUptime[0].Count)
}

func TestValidatorUptimeService_ListUptimeByOperAddresses(t *testing.T) {

	validatorAddresses := []string{"omnivaloper195tf7tuvxagrsq7seh8579qrjvxc6j9wws29sg", "omnivaloper1ax6d7u25ts4yuv03d8xlcymwyy5uyhcru3tkuv"}
	validatorUptime, err := ValidatorUptimeService.ListRecentUptimeByOperAddresses(context.Background(), validatorAddresses, 500)
	if err != nil {
		t.Fatalf("failed to get validator uptime: %v", err)
	}
	fmt.Println(validatorUptime[0].Count)
}

func TestBech32ToHex(t *testing.T) {

	addr := "omnivalcons1vcumflh0a6wm0e7aa8zkjcdge6zeydq3h60y0w"
	hexAddr, err := ValidatorUptimeService.Bech32ToHex(addr, "omnivalcons")
	if err != nil {
		t.Fatalf("failed to get validator uptime: %v", err)
	}
	fmt.Println(strings.ToUpper(hexAddr))
}
