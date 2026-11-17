package favorites

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	// user пользователь из переменных окружения, он работает со своим избранным
	user models.LoginOkResponse
	// seller второй пользователь, владелец объявлений: свои объявления в избранное добавить нельзя
	seller models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения")
	s.user = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))

	base.Precondition("Регистрация второго пользователя, который будет владельцем объявлений")
	s.seller = auth.RegisterRandomUser(s.T())
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// createSellerAdvertisement создаёт объявление второго пользователя и регистрирует его удаление после теста.
// Вызывать нужно на уровне самого теста, а не внутри Step.
func (s *TestSuite) createSellerAdvertisement() advModels.AdvertisementResponse {
	advertisement := advertisements.CreateAdvertisement(s.T(), s.seller.Token, advertisements.NewAdvertisementRequest(1))
	base.DeleteAdvertisementAfterTest(s.T(), s.seller.Token, advertisement.Id)

	return advertisement
}
