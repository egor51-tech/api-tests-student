package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/constants/testdata"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

// Login выполняет авторизацию пользователя и возвращает токен с данными пользователя, ожидает 200
func Login(t *testing.T, username, password string) models.LoginOkResponse {
	require.NotEmpty(t, username, "Отсутствует логин пользователя")
	require.NotEmpty(t, password, "Отсутствует пароль пользователя")

	request := models.LoginRequest{
		Email:    username,
		Password: password,
	}

	loginRawResponse := auth.HttpPostAuthLogin(t, request)
	defer loginRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(loginRawResponse.Body)
	require.NoError(t, errResponse)

	loginResponse := models.LoginOkResponse{}
	err := json.Unmarshal(respArray, &loginResponse)
	require.NoError(t, err)

	return loginResponse
}

// Register регистрирует нового пользователя и возвращает токен с его данными, ожидает 201
func Register(t *testing.T, request models.RegisterRequest) models.LoginOkResponse {
	t.Helper()

	registerRawResponse := auth.HttpPostAuthRegister(t, request)
	body := response.ReadBodyWithStatus(t, registerRawResponse, http.StatusCreated)

	registerResponse := response.Decode[models.LoginOkResponse](t, body)
	require.NotEmpty(t, registerResponse.Token, "В ответе на регистрацию отсутствует токен")
	require.NotEmpty(t, registerResponse.User.Id, "В ответе на регистрацию отсутствует id пользователя")

	return registerResponse
}

// RegisterRandomUser регистрирует пользователя со случайными email и паролем;
// используется там, где нужен второй пользователь, например владелец чужого объявления
func RegisterRandomUser(t *testing.T) models.LoginOkResponse {
	t.Helper()

	email := fmt.Sprintf("qa_%s@%s", strings.ToLower(utils.RandomString(testdata.RandomSuffixLength)), testdata.UserEmailDomain)

	return Register(t, models.RegisterRequest{
		Email:     email,
		Password:  utils.RandomPassword(testdata.UserPasswordLength),
		FirstName: testdata.UserFirstName,
	})
}
