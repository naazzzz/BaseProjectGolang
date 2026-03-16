package authmdl

import (
	"BaseProjectGolang/internal/config"
	common "BaseProjectGolang/internal/constant"
	cryptoUtil "BaseProjectGolang/pkg/crypto"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dromara/carbon/v2"
	"github.com/gofiber/fiber/v3"
	"github.com/rotisserie/eris"
)

const (
	TimeDifferenceForAuthorizeRequest = 10
)

type QuerySignatureWithTime struct {
	cfg *config.Config
}

func NewQuerySignatureWithTime(
	cfg *config.Config,
) *QuerySignatureWithTime {
	return &QuerySignatureWithTime{
		cfg: cfg,
	}
}

// AuthorizeRequest Для авторизации использовался https://en.wikipedia.org/wiki/HMAC c hash функцией SHA256
// Не отправлять в Query значения с одинаковыми ключами - php их перезаписывает, go - формирует массив,
// также не отправлять вложенные структуры типа a[b][c] для го это отдельный query,
// либо ключи эл-тов в языках будут отличаться
func (auth *QuerySignatureWithTime) AuthorizeRequest(ctx fiber.Ctx) error {
	signatureTimestamp, err := strconv.Atoi(ctx.Query(common.SignatureTimestampParamsKey))
	if err != nil {
		return eris.Wrap(fiber.NewError(http.StatusUnauthorized, err.Error()), http.StatusText(fiber.StatusUnauthorized))
	}

	signature := ctx.Query(common.SignatureParamsKey)
	if signature == "" {
		return eris.Wrap(fiber.NewError(http.StatusUnauthorized, "Empty signature"), http.StatusText(fiber.StatusUnauthorized))
	}

	query, err := url.ParseQuery(string(ctx.Request().URI().QueryString()))
	if err != nil {
		return eris.Wrap(fiber.NewError(http.StatusUnauthorized, err.Error()), http.StatusText(fiber.StatusUnauthorized))
	}

	jsonString, _ := FormPublicKey(query, ctx.Body(), signatureTimestamp)
	publicKey := cryptoUtil.GetMD5Hash(jsonString)
	valid := cryptoUtil.ValidateSignatureHmacSha256(signature, publicKey, auth.cfg.Secure.AuthPrivateKey)
	requestTime := carbon.CreateFromTimestamp(int64(signatureTimestamp))

	if valid && requestTime.Between(
		carbon.Now().SubSeconds(TimeDifferenceForAuthorizeRequest),
		carbon.Now().AddSeconds(TimeDifferenceForAuthorizeRequest),
	) {
		return ctx.Next()
	}

	return eris.Wrap(fiber.NewError(http.StatusUnauthorized, ""), http.StatusText(fiber.StatusUnauthorized))
}

func FormPublicKey(
	queries map[string][]string,
	body []byte,
	signatureTimestamp int,
) (string, error) {
	var jsonBody map[string]interface{}

	err := json.Unmarshal(body, &jsonBody)
	if err != nil {
		jsonBody = map[string]interface{}{}
	}

	jsonBody["signature_timestamp"] = signatureTimestamp

	for key, value := range queries {
		if key == "signature_timestamp" || key == "signature" {
			continue
		}

		if 1 < len(value) {
			jsonBody[key] = value
		} else {
			jsonBody[key] = value[0]
		}
	}

	result, err := json.Marshal(jsonBody)
	if err != nil {
		return "", err
	}

	return string(result), err
}
