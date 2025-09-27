package utils_test

import (
	"fmt"
	"net/http"
	"testing"

	"google.golang.org/api/oauth2/v2"
)

func TestJwt(t *testing.T) {
	s := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjhjMjdkYjRkMTNmNTRlNjU3ZDI2NWI0NTExMDA4MGI0ODhlYjQzOGEiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJuYmYiOjE2NzEzNzE5OTgsImF1ZCI6Ijk0NjU2MjY4NjAwNS12NzYzOGgxaTJwM2t0ZzNybmUyb21hb3A4YmMydWZsYS5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsInN1YiI6IjExMzM5MjYwNDgyMzM2NDE1NzgxMCIsImVtYWlsIjoiNDIwNjQyMjA3QHFxLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhenAiOiI5NDY1NjI2ODYwMDUtdjc2MzhoMWkycDNrdGczcm5lMm9tYW9wOGJjMnVmbGEuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJuYW1lIjoiNDIwNjQyMjA3QHFxLmNvbSIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS0vQUQ1LVdDbGthNFc2UXBZRzVpX0tIY3haMkNBTk11VTNidGxOVXczenFMVGg9czk2LWMiLCJpYXQiOjE2NzEzNzIyOTgsImV4cCI6MTY3MTM3NTg5OCwianRpIjoiYjk1Njk1NTNkOTc3NWFmZDM2ODU0YjY0YTE4ZjRjOTU0ZDlhMDFmZCJ9.cTM9ujPvmfCzj0Ik8x8uZByDH5QgomdWtc6gP3cz4WqZtTliB5DH5t48gfXOu6GbDXy9AsoWSYfj70k2mGinauEmHYXgVH3o7YQTF8Trvky2CjeEyeiBKLFdyCxBQwa_-lBDAF7t8ACt0W9oJFX9LKiDxkZvgJi4I8MGPDkph19QXvIYCtdkA3qO1n7t-EYEf2Maun7DYfqEDfCKoQvYdHqqmIf_yvRGq6Yq4XP7wSx4W8ea4-cuClx7iq9tROZQxAHyA4wrKe-VhMHRKH2vqC0fXgQwfhURGh2JMRKsOfG2S6JQYqUP9KvF3UEs4ZDFX7J3HlQqFRgx_1uiAEDcOQ"
	client, err := oauth2.New(&http.Client{})
	if err != nil {
		t.Error(err)
	}
	tokenInfoCall := client.Tokeninfo()
	tokenInfoCall.IdToken(s)
	info, err := tokenInfoCall.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("email:", info.Email)
	fmt.Println("userID:", info.UserId)
	fmt.Println("aud:", info.Audience)

	j, _ := info.MarshalJSON()
	fmt.Println("json:", j)
	fmt.Println("info:", info)
	t.Fail()

}
