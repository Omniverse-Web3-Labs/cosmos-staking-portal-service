package login

import (
	"app/internal/http/resp"
	"app/internal/service"
	"strconv"
	"strings"

	"encoding/base64"
	"encoding/hex"
	"fmt"

	arbitraryverify "github.com/MyriadFlow/cosmos-wallet/sign-auth/pkg/cosmos_blockchain/arbitrary_verify"
	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/gin-gonic/gin"
)

type CosmosWalletRequest struct {
	Address    string `json:"address" binding:"required"`
	Pubkey     string `json:"pubkey" binding:"required"`
	Signature  string `json:"sig" binding:"required"`
	WalletUUID string `json:"walletUUID"`
	WalletName string `json:"walletName"`
}

func verifyCosmosSignature(address, pubkeyHexString, msg, sigBase64 string) (bool, error) {
	// decode
	pubkeyBytes, err := hex.DecodeString(pubkeyHexString)
	if err != nil {
		fmt.Println("decode pubkey error:", err)
		return false, err
	}

	pk := &secp256k1.PubKey{Key: pubkeyBytes}

	// decrypt
	sigBytes, err := base64.StdEncoding.DecodeString(sigBase64)
	if err != nil {
		fmt.Println("decode sig error:", err)
		return false, err
	}

	// verify signature
	ok, err := arbitraryverify.VerifyArbitraryMsg(
		address,
		msg,
		sigBytes,
		*pk,
	)

	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	return true, nil
}

func CosmosWalletAction(c *gin.Context) {
	ctx := c.Request.Context()
	var request CosmosWalletRequest
	var data any
	err := func() error {
		if err := c.ShouldBindJSON(&request); err != nil {
			return err
		}

		address := strings.ToLower(request.Address)

		// it should be the same as the challenge fetched before
		challenge, err := service.UserService.GetChallenge(ctx, address)
		if err != nil {
			return err
		}
		message := fmt.Sprintf("Sign in to Cosmos Staking!\nChallenge:%s", challenge)

		verified, err := verifyCosmosSignature(address, request.Pubkey, message, request.Signature)
		if err != nil {
			return err
		}

		if !verified {
			return fmt.Errorf("failed to verify address:%s", address)
		}

		// try to create user
		user, err := service.UserService.TryCreateUser(ctx, request.Address, request.WalletUUID, request.WalletName)

		if err != nil {
			return err
		}
		authInfo, err := service.AuthService.Login(ctx, strconv.FormatInt(user.ID, 10), address)
		if err != nil {
			return err
		}

		if err := service.UserService.UpdateChallenge(ctx, user.ID, address); err != nil {
			return err
		}
		data = authInfo
		return nil
	}()
	resp.Finish(c, err, data)
}
