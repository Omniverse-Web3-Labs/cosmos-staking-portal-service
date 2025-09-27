package login

import (
	"fmt"
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const address = "omni195tf7tuvxagrsq7seh8579qrjvxc6j9wgnpz86"
const pubkey = "020590a43e1008d407a739000048680484675084210b00d1cf8781a44f8d12cd8c"
const message = "Sign in to Cosmos Staking!\nChallenge:0000-0000-0000-0000"
const signature = "RU8TFCP8Nim2OHpSnXtmxSyQXwM517LpbSy9zri+Tt1vC08kGULdedroUb1Svylczu6YcVtLjWmnPO/MFTQCzQ=="

func TestMain(m *testing.M) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("omni", "omnipub")
	config.Seal()

	exitCode := m.Run()

	fmt.Println("cleaning...")

	os.Exit(exitCode)
}

func TestVerifySignatureFailWithPubkeyError(t *testing.T) {
	pk := []byte(pubkey)
	pk[0] = '-'
	_, err := verifyCosmosSignature(address, string(pk), message, signature)
	if err == nil {
		t.Errorf("Expected an error but got nil")
	}
}

func TestVerifySignatureFailWithSignatureError(t *testing.T) {
	sig := []byte(signature)
	sig[0] = '-'
	_, err := verifyCosmosSignature(address, pubkey, message, string(sig))
	if err == nil {
		t.Errorf("Expected an error but got nil")
	}
}

func TestVerifySignatureFailWithVerifyError(t *testing.T) {
	msg := []byte(message)
	msg[0] = 'G'
	pass, err := verifyCosmosSignature(address, pubkey, string(msg), signature)
	if err != nil {
		t.Errorf("Not expected error")
	}

	if pass {
		t.Errorf("The address should not match")
	}
}

func TestVerifySignatureFailWithAddressError(t *testing.T) {
	pass, err := verifyCosmosSignature("address", pubkey, message, signature)
	if err != nil {
		t.Errorf("Not expected error")
	}

	if pass {
		t.Errorf("The address should not match")
	}
}

func TestVerifySignaturePass(t *testing.T) {
	pass, err := verifyCosmosSignature(address, pubkey, message, signature)
	if err != nil {
		t.Errorf("Not expected error")
	}

	if !pass {
		t.Errorf("The address should match")
	}
}
