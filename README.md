# API-тесты для проекта https://testboard.avito.com/

## Полезные ссылки

1. [Swagger-спецификация](https://testboard.avito.com/swagger/index.html)
2. [Техническое задание](https://docs.google.com/document/d/1IZ9SxEcePoQT6lyEmukycp1qNok4HZtok73UQnPwKSA/edit?usp=sharing)
3. [Библиотека для написания API-тестов](https://github.com/steinfletcher/apitest)
4. [Мощная библиотека анализа JSON для Go](https://github.com/tidwall/gjson)

## Запуск тестов

1. Установите любое IDE для разработки: [Jetbrains GoLand](https://www.jetbrains.com/ru-ru/go/) (Платное или по студенческой лицензии), [VS Code](https://code.visualstudio.com/) (бесплатное).
2. Установите [Git](https://git-scm.com/install/) на компьютер
3. [Установите](https://go.dev/doc/install) Go версии 1.25.1 или выше
4. Склонируйте репозиторий по HTTPS (репозиторий публичный, ключи не нужны):

```bash
  git clone https://github.com/<ваш-логин>/api-tests-student.git
  cd api-tests-student
```

5. Подтяните необходимые зависимости через команду в корне проекта:

```bash
  go mod tidy
```

6. В файле `.env.override` (создайте его в корне, если нет) выставьте логин и пароль от [TestBoard](https://testboard.avito.com/):

```
TEST_LOGIN=your_login
TEST_PASSWORD=your_password
```

7. Для запуска всех тестов достаточно выполнить команду в корне проекта:

```bash
  go test -count=1 -v ./...
```

8. Отдельный пакет с тестами можно запустить командой:

```bash
  go test -count=1 -v ./tests/scenarios/createAdvertisement/
  go test -count=1 -v ./tests/scenarios/favorites/
  go test -count=1 -v ./tests/scenarios/myAdvertisement/
```

Флаг `-count=1` отключает кеш результатов Go, иначе повторный запуск без изменений кода вернёт старый результат,
а не выполнит тесты заново.

## Структура проекта

Тесты построены по слоям `Client Layer → Manager Layer → Scenario`:

```
internal/
├── client/http/         # Client Layer: HTTP-запросы к ручкам (auth, advertisements, favorites, myAdvertisements)
├── managers/            # Manager Layer: вызов клиента, проверка статус-кода, разбор ответа в модели
├── constants/           # пути ручек, тестовые данные, тексты ошибок
├── helpers/             # api-runner (обёртка над apitest) и разбор HTTP-ответов
└── utils/               # загрузка env, логирование, генерация случайных строк и PNG-картинок
tests/
├── base.go              # общие действия сьюта, шаги с логированием, очистка данных после теста
└── scenarios/           # Scenario: сами тесты, один пакет на один сценарий
    ├── createAdvertisement/   # Часть 1. Добавление объявления
    ├── favorites/             # Часть 2. Работа с избранным
    └── myAdvertisement/       # пример из шаблона
```

## Что покрыто

### Часть 1. Добавление объявления (`createAdvertisement`)

| Тест | Что проверяет |
|------|---------------|
| `TestCreateAdvertisementPositive` | Объявление со всеми параметрами и тремя фотографиями: `GET /advertisement`, поиск по полному названию, `GET /advertisements/{id}/photos` |
| `TestCreateAdvertisementInvalidParams` | 400 и сообщение об ошибке: пустое название, название длиннее 255 символов, пустое описание, отрицательная цена, нет фотографий, фотографий больше трёх |
| `TestCreateAdvertisementUnauthorized` | 401 без токена и с неверным токеном |

### Часть 2. Работа с избранным (`favorites`)

| Тест | Что проверяет |
|------|---------------|
| `TestFavoritesFullCycle` | Добавление, появление в `GET /favorites`, удаление, исчезновение из списка |
| `TestDeleteAdvertisementRemovesItFromFavorites` | После `DELETE /advertisement` объявление пропадает из избранного |
| `TestAddNonExistentAdvertisement`, `TestRemoveNonExistentAdvertisement` | 404 при добавлении и удалении несуществующего `ad_id` |
| `TestAddOwnAdvertisement` | 403 при добавлении собственного объявления |
| `TestAddAdvertisementTwice` | 409 при повторном добавлении |
| `TestAddWithoutAdvertisementId` | 400 при пустом `advertisement_id` |
| `TestAddOverFavoritesLimit` | 400 при добавлении 51-го объявления (лимит избранного - 50) |
| `TestFavoritesUnauthorized` | 401 для `GET`, `POST` и `DELETE` без токена и с неверным токеном |

Ответ 500 из swagger на практике не воспроизводится, поэтому отдельного теста на него нет.

## Особенности

- Свои объявления в избранное добавить нельзя, поэтому для сценариев избранного в `SetupSuite` регистрируется второй
  пользователь со случайным email. Он создаёт объявления, а пользователь из `.env.override` добавляет их в избранное.
  Регистрироваться заново приходится на каждый запуск, потому что удалить пользователя может только администратор.
- Каждый тест сам создаёт нужные данные с уникальными названиями и удаляет их в `t.Cleanup`, поэтому тесты можно
  запускать повторно и в любом порядке.
- Тест `TestAddOverFavoritesLimit` создаёт 51 объявление, поэтому выполняется заметно дольше остальных.
- Картинки для загрузки генерируются в памяти, файлы в репозитории не нужны.

## Логирование

Тесты пишут в консоль подготовку, шаги и результаты проверок с метками `[Precondition]`, `[Step]`, `[Case]`,
`[Result]` и `[Cleanup]`.

Для включения логирования request/response в консоль: выставьте переменную окружения `DEBUG=true` в
`.env.override`

## Запуск линтера

1. Установить приложение [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)
2. Запуск линтера (из корня репозитория):

```bash
  golangci-lint run
```
