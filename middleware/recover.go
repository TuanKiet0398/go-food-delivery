package middleware

import (
	"fmt"

	"food-delivery/common"
	"food-delivery/component/appctx"

	"github.com/gin-gonic/gin"
)

func Recover(sc appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Content-Type", "application/json")

				if appErr, ok := err.(*common.AppError); ok {
					c.AbortWithStatusJSON(appErr.StatusCode, appErr)
					return
				}

				rootErr, ok := err.(error)
				if !ok {
					rootErr = fmt.Errorf("%v", err)
				}

				appErr := common.ErrInternal(rootErr)
				c.AbortWithStatusJSON(appErr.StatusCode, appErr)
			}
		}()

		c.Next()
	}
}