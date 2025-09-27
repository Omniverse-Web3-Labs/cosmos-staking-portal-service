package login

import (
	"app/internal/http/resp"
	"app/internal/service"
	"strconv"
	"strings"

	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

type WalletRequest struct {
	Address    string `json:"address" binding:"required"`
	Signature  string `json:"sig" binding:"required"`
	WalletUUID string `json:"walletUUID"`
	WalletName string `json:"walletName"`
}

func verifySignature(address string, message string, signatureHex string) (bool, error) {
	signature, err := hex.DecodeString(signatureHex[2:])
	if err != nil {
		return false, err
	}

	if len(signature) != 65 {
		return false, errors.New("signature error")
	}

	if signature[64] >= 27 {
		signature[64] -= 27
	}

	hash := crypto.Keccak256Hash([]byte(message))

	publicKey, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		return false, err
	}

	recovered := crypto.PubkeyToAddress(*publicKey).Hex()

	if strings.ToLower(address) != strings.ToLower(recovered) {
		return false, nil
	} else {
		return true, nil
	}
}

func WalletAction(c *gin.Context) {
	ctx := c.Request.Context()
	var request WalletRequest
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
		prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))

		verified, err := verifySignature(address, prefix+message, request.Signature)
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
