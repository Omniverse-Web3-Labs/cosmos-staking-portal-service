package keybase

import (
	"app/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
)

type keybaseStatus struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

type keybaseKeyResponse struct {
	Status keybaseStatus `json:"status"`
	Keys   []keybaseKey  `json:"keys"`
}

type keybaseKey struct {
	KeyFingerprint string `json:"key_fingerprint"`
}

type keybaseUserResponse struct {
	Status keybaseStatus `json:"status"`
	Them   []keybaseUser `json:"them"`
}

type keybaseUser struct {
	Basics struct {
		Username string `json:"username"`
	} `json:"basics"`
	Pictures struct {
		Primary struct {
			URL string `json:"url"`
		} `json:"primary"`
	} `json:"pictures"`
}

type keybase struct {
}

var Keybase keybase

func (k *keybase) GetAvatar(ctx context.Context, identity string) (string, error) {
	url := fmt.Sprintf("https://keybase.io/_/api/1.0/key/fetch.json?pgp_key_ids=%s", identity)

	data, err := utils.HttpUtil.SendRequestWithContext(ctx, url, "GET", nil, nil)
	if err != nil {
		return "", err
	}

	var keyResponse keybaseKeyResponse
	err = json.Unmarshal(data, &keyResponse)
	if err != nil {
		return "", err
	}

	if keyResponse.Status.Code != 0 {
		return "", fmt.Errorf("query keybase keys of %s faild: %s", identity, keyResponse.Status.Name)
	}

	if len(keyResponse.Keys) == 0 {
		return "", fmt.Errorf("there is no key of %s", identity)
	}

	url = fmt.Sprintf("https://keybase.io/_/api/1.0/user/lookup.json?key_fingerprint=%s", keyResponse.Keys[0].KeyFingerprint)

	data, err = utils.HttpUtil.SendRequestWithContext(ctx, url, "GET", nil, nil)
	if err != nil {
		return "", err
	}

	var userResponse keybaseUserResponse
	err = json.Unmarshal(data, &userResponse)
	if err != nil {
		return "", err
	}

	if userResponse.Status.Code != 0 {
		return "", fmt.Errorf("query keybase user of %s faild: %s", keyResponse.Keys[0].KeyFingerprint, userResponse.Status.Name)
	}

	if len(userResponse.Them) == 0 {
		return "", fmt.Errorf("there is no user of %s", keyResponse.Keys[0].KeyFingerprint)
	}

	return userResponse.Them[0].Pictures.Primary.URL, nil
}
