//  *@createTime    2022/3/21 23:36
//  *@author        hay&object
//  *@version       v1.0.0

package api

import (
	"errors"
	"fmt"
	"github.com/haylove/small_bank/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationTypeBearer = "bearer"
	authorizationHeaderKey  = "authorization"
	authorizationPayloadKey = "authorization_payload_key"
)

// authMiddleware returns a middleware with authorization
func authMiddleware(generator token.Generator) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// get authorizationHeader if existed, if no header ,that means 401
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		// check authorizationHeader format , should be like "Type Token"
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authentication header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		// check authorizationHeader type
		authorizationType := strings.ToLower(fields[0])
		accessToken := fields[1]
		if !isSupportedType(authorizationType) {
			err := fmt.Errorf("unsupported authentication type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		// verify token
		payload, err := generator.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

// isSupportedType returns the authorizationType is it supported
func isSupportedType(authorizationType string) bool {
	switch authorizationType {
	case authorizationTypeBearer:
		return true
	}
	return false
}
