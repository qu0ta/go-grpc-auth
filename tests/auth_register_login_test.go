package tests

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/qu0ta/go-grpc-auth/tests/suite"
	authv1 "github.com/qu0ta/pet-proto/gen/go/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	appId          = 1
	appSecret      = "secret1"
	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	t.Run("HappyPath", func(t *testing.T) {
		email := gofakeit.Email()
		password := fakePassword()

		respReg, err := st.AuthClient.Register(ctx, &authv1.RegisterRequest{
			Email:    email,
			Password: password,
			AppId:    appId,
		})

		require.NoError(t, err)
		assert.NotEmpty(t, respReg.GetUserId())

		respLogin, err := st.AuthClient.Login(ctx, &authv1.LoginRequest{
			Email:    email,
			Password: password,
		})

		loginTime := time.Now()
		require.NoError(t, err)

		token := respLogin.GetToken()
		assert.NotEmpty(t, token)

		tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(appSecret), nil
		})
		require.NoError(t, err)

		claims, ok := tokenParsed.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, email, claims["email"].(string))
		assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
		assert.Equal(t, appId, int(claims["app_id"].(float64)))

		const deltaSec = 1
		assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSec)
	})

	t.Run("WrongPassword", func(t *testing.T) {
		email := gofakeit.Email()
		password := fakePassword()

		_, err := st.AuthClient.Register(ctx, &authv1.RegisterRequest{
			Email:    email,
			Password: password,
			AppId:    appId,
		})

		require.NoError(t, err)

		_, err = st.AuthClient.Login(ctx, &authv1.LoginRequest{
			Email:    email,
			Password: "wrong-password",
		})

		assert.Error(t, err)

	})

	t.Run("UnauthorizedEmail", func(t *testing.T) {
		_, err := st.AuthClient.Login(ctx, &authv1.LoginRequest{
			Email:    gofakeit.Email(),
			Password: fakePassword(),
		})

		assert.Error(t, err)
	})

	t.Run("BlankEmailRegister", func(t *testing.T) {
		_, err := st.AuthClient.Register(ctx, &authv1.RegisterRequest{
			Email:    "",
			Password: fakePassword(),
			AppId:    appId,
		})

		assert.Error(t, err)
	})

	t.Run("BlankEmailLogin", func(t *testing.T) {
		_, err := st.AuthClient.Login(ctx, &authv1.LoginRequest{
			Email:    "",
			Password: fakePassword(),
		})

		assert.Error(t, err)
	})

	t.Run("BlankPasswordRegister", func(t *testing.T) {
		_, err := st.AuthClient.Register(ctx, &authv1.RegisterRequest{
			Email:    gofakeit.Email(),
			Password: "",
			AppId:    appId,
		})

		assert.Error(t, err)
	})

	t.Run("BlankPasswordLogin", func(t *testing.T) {
		_, err := st.AuthClient.Login(ctx, &authv1.LoginRequest{
			Email:    gofakeit.Email(),
			Password: "",
		})

		assert.Error(t, err)
	})

	t.Run("DoubleEmailRegistration", func(t *testing.T) {
		email := gofakeit.Email()
		password := fakePassword()
		_, err := st.AuthClient.Register(ctx, &authv1.RegisterRequest{
			Email:    email,
			Password: password,
			AppId:    appId,
		})

		require.NoError(t, err)

		_, err = st.AuthClient.Register(ctx, &authv1.RegisterRequest{
			Email:    email,
			Password: password,
			AppId:    appId,
		})

		assert.Error(t, err)
	})

}
func fakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
