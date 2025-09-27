package service

import (
	"app/internal/db/postgres/model"
	"app/internal/db/postgres/query"
	"app/pkg/utils"
	"log/slog"
	"strings"

	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"gorm.io/gorm"
	// "app/pkg/utils"
	// "strings"
	// "github.com/gin-gonic/gin"
)

var UserService userService

type userService struct {
}

func generateChallenge() (string, error) {
	// 创建一个 8 字节的随机字节切片
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// 将字节数组转换为十六进制字符串
	hexString := hex.EncodeToString(randomBytes)

	// 返回格式化后的字符串，加入分隔符 '-'
	return fmt.Sprintf("%s-%s-%s-%s", hexString[:4], hexString[4:8], hexString[8:12], hexString[12:]), nil
}

func (service *userService) GetChallenge(ctx context.Context, address string) (string, error) {
	users, err := query.User.WithContext(ctx).Where(query.User.Address.Eq(address)).Find()
	if err != nil {
		return "", err
	}

	if len(users) == 0 {
		return "0000-0000-0000-0000", nil
	} else {
		return users[0].Challenge, nil
	}
}

func (service *userService) UpdateChallenge(ctx context.Context, uid int64, address string) error {
	challenge, err := generateChallenge()
	if err != nil {
		return err
	}

	utils.Logger.Info(ctx, "user", "UpdateChallenge", "update challenge", slog.Attr{
		Key:   "uid",
		Value: slog.Int64Value(uid),
	}, slog.Attr{
		Key:   "challenge",
		Value: slog.StringValue(challenge),
	})

	_, errUpdate := query.User.WithContext(ctx).Where(query.User.Address.Eq(address)).UpdateSimple(query.User.Challenge.Value(challenge))

	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (service *userService) TryCreateUser(ctx context.Context, address, walletUUID, walletName string) (*model.User, error) {
	lowerCase := strings.ToLower(address)
	user, err := query.User.WithContext(ctx).Where(query.User.Address.Eq(lowerCase)).Take()
	if err == gorm.ErrRecordNotFound {
		if lowerCase != address {
			user, err = query.User.WithContext(ctx).Where(query.User.Address.Eq(address)).Take()
			if err == nil {
				utils.Logger.Info(ctx, "user", "TryCreateUser", "found user with not lowercase address", slog.Attr{
					Key:   "uid",
					Value: slog.Int64Value(user.ID),
				}, slog.Attr{
					Key:   "address",
					Value: slog.StringValue(address),
				})
				return user, nil
			} else if err != gorm.ErrRecordNotFound {
				return nil, err
			}
		}
		createdUser := model.User{
			Address:    lowerCase,
			Challenge:  "0000-0000-0000-0000",
			WalletUUID: walletUUID,
			WalletName: walletName,
		}

		utils.Logger.Info(ctx, "user", "TryCreateUser", "create user", slog.Attr{
			Key:   "address",
			Value: slog.StringValue(lowerCase),
		})

		err := query.User.WithContext(ctx).Create(&createdUser)
		if err != nil {
			return nil, err
		}
		return &createdUser, nil
	} else if err != nil {
		return nil, err
	} else {
		if user.WalletUUID == "" || user.WalletName == "" {
			query.User.WithContext(ctx).
				Where(query.User.Address.Eq(lowerCase)).
				UpdateSimple(query.User.WalletUUID.Value(walletUUID), query.User.WalletName.Value(walletName))
		}
	}
	return user, nil
}

func (server *userService) GetUserAddress(ctx context.Context, uid int64) (string, error) {
	user, err := query.User.WithContext(ctx).Where(query.User.ID.Eq(uid)).Find()

	if err != nil {
		return "", err
	}

	if len(user) == 0 {
		return "", fmt.Errorf("user not found: %d", uid)
	}

	return user[0].Address, nil
}
