// authorization.go kee > 2019/12/15

package routes

import (
	"encoding/json"
	"koobeton/app"
	c "koobeton/controllers"
	m "koobeton/models"
	"strconv"
	"time"
)

func authorization(ctx app.Context) {
	header := ctx.Request().Header
	// 签名到期时间 - 客户端控制凭证到期
	expireStr := header.Get("X-Expire")
	// 授权token
	auth := header.Get("Authorization")
	pk := header.Get("X-PaddingKey")
	xToken := header.Get("X-Token")

	// 校验签名到期时间
	if expire, _ := strconv.Atoi(expireStr); expire > 0 && time.Now().Unix() > int64(expire) {
		authFail(ctx, c.ResError(401, "token out of date"))
		return
	}

	if auth != "" {
		ckAuthorization(ctx, auth, expireStr, pk)
	} else if xToken != "" {
		ckToken(ctx, xToken)
	}
}

// auth = accessKey + SHA1(AccessKey:SecretKey + Pk (Len(Pk) >= 0))
func ckAuthorization(ctx app.Context, auth, expireStr, pk string) {
	if auth != "" && len(auth) > 32 {
		accessKey := auth[:32]
		token := auth[32+1:]
		// 校验token
		uToken := &m.UserToken{AccessKey: accessKey}
		user, err := uToken.CheckToken(token, expireStr, pk)
		if err != nil {
			authFail(ctx, c.ResError(401, err.Error()))
			return
		}
		uCode, _ := json.Marshal(user)
		// 写入请求
		ctx.Params().Set("cUser", string(uCode))
	} else {
		authFail(ctx, c.ResError(401, "token not specified"))
		return
	}
	ctx.Next()
}

func ckToken(ctx app.Context, token string) {
	user := &m.User{}
	data, err := user.VerifyToken(token)
	if err != nil {
		authFail(ctx, c.ResError(401, err.Error()))
		return
	}
	userToken := m.UserToken{}
	claims := new(app.Json).Conver(data.Get("Claims"))
	userToken.AccessKey = claims.Get("accessKey").(string)
	user, err = userToken.GetUser()
	if err != nil {
		authFail(ctx, c.ResError(401, "bad token - "+err.Error()))
		return
	}
	uCode, _ := json.Marshal(user)
	ctx.Params().Set("cUser", string(uCode))
	ctx.Next()
}

func authFail(ctx app.Context, result c.Result) {
	ctx.StatusCode(result.Code)
	ctx.JSON(result)
}
