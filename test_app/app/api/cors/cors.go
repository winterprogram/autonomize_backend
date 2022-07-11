package cors

import (
	"strconv"
	"strings"

	"test/test_app/app/service/correlation"
	"test/test_app/app/service/logger"

	"github.com/gin-gonic/gin"
)

// Options is struct that defined cors properties
type Options struct {
	AllowedOrigins   []string
	AllowedHeaders   []string
	AllowMaxAge      int
	AllowedMethods   []string
	ExposedHeaders   []string
	AllowCredentials bool
	Debug            bool
}

type CorsHandler struct {
	allowedOriginsAll bool
	allowedOrigins    []string
	allowedHeadersAll bool
	allowedHeaders    []string
	allowedMethods    []string
	exposedHeaders    []string
	allowCredentials  bool
	maxAge            int
}

var defaultOptions = &Options{
	AllowedOrigins: []string{"*"},
	AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE, PUT"},
	AllowedHeaders: []string{"Origin", "Accept", "Content-Type"},
}

func DefaultHandler() *CorsHandler {
	return NewCorsHandler(*defaultOptions)
}

func NewCorsHandler(options Options) *CorsHandler {
	cors := &CorsHandler{
		allowedOrigins:   options.AllowedOrigins,
		allowedHeaders:   options.AllowedHeaders,
		allowCredentials: options.AllowCredentials,
		allowedMethods:   options.AllowedMethods,
		exposedHeaders:   options.ExposedHeaders,
		maxAge:           options.AllowMaxAge,
	}

	if len(cors.allowedOrigins) == 0 {
		cors.allowedOrigins = defaultOptions.AllowedOrigins
		cors.allowedOriginsAll = true
	} else {
		for _, v := range options.AllowedOrigins {
			if v == "*" {
				cors.allowedOrigins = defaultOptions.AllowedOrigins
				cors.allowedOriginsAll = true
				break
			}
		}
	}
	if len(cors.allowedHeaders) == 0 {
		cors.allowedHeaders = defaultOptions.AllowedHeaders
		cors.allowedHeadersAll = true
	} else {
		for _, v := range options.AllowedHeaders {
			if v == "*" {
				cors.allowedHeadersAll = true
				break
			}
		}
	}
	if len(cors.allowedMethods) == 0 {
		cors.allowedMethods = defaultOptions.AllowedMethods
	}
	return cors
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func (c *CorsHandler) CorsMiddleware(innerHandler gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if string(ctx.Request.Method) == "OPTIONS" {
			c.handlePreflight(ctx)
			ctx.Status(200)
		} else {
			c.handleActual(ctx)
			innerHandler(ctx)
		}
	}
}

func stringConcat(a []string) string {
	return strings.Join(a, ",")
}

func (c *CorsHandler) handlePreflight(ct *gin.Context) {
	originHeader := string(ct.GetHeader("Origin"))
	ctx := correlation.WithReqContext(ct)
	log := logger.Logger(ctx)
	if len(originHeader) == 0 || !c.isAllowedOrigin(originHeader) {
		log.Info("Origin " + originHeader + " is not in " + stringConcat(c.allowedOrigins))
		return
	}
	method := string(ct.GetHeader("Access-Control-Request-Method"))
	if !c.isAllowedMethod(method) {
		log.Info("Method " + method + " is not in" + stringConcat(c.allowedMethods))
		return
	}
	headers := []string{}
	if len(ct.GetHeader("Access-Control-Request-Headers")) > 0 {
		headers = strings.Split(string(ct.GetHeader("Access-Control-Request-Headers")), ",")
	}
	if !c.areHeadersAllowed(headers) {
		log.Info("Headers " + stringConcat(headers) + " is not in" + stringConcat(c.allowedHeaders))
		return
	}

	ct.Set("Access-Control-Allow-Origin", originHeader)
	ct.Set("Access-Control-Allow-Methods", method)
	if len(headers) > 0 {
		ct.Set("Access-Control-Allow-Headers", strings.Join(headers, ", "))
	}
	if c.allowCredentials {
		ct.Set("Access-Control-Allow-Credentials", "true")
	}
	if c.maxAge > 0 {
		ct.Set("Access-Control-Max-Age", strconv.Itoa(c.maxAge))
	}
}

func (c *CorsHandler) handleActual(ct *gin.Context) {
	originHeader := string(ct.GetHeader("Origin"))
	ctx := correlation.WithReqContext(ct)
	log := logger.Logger(ctx)

	if len(originHeader) == 0 || !c.isAllowedOrigin(originHeader) {
		log.Info("Origin " + originHeader + " is not in " + stringConcat(c.allowedOrigins))
		return
	}
	ct.Set("Access-Control-Allow-Origin", originHeader)
	if len(c.exposedHeaders) > 0 {
		ct.Set("Access-Control-Expose-Headers", strings.Join(c.exposedHeaders, ", "))
	}
	if c.allowCredentials {
		ct.Set("Access-Control-Allow-Credentials", "true")
	}
}

func (c *CorsHandler) isAllowedOrigin(originHeader string) bool {
	if c.allowedOriginsAll {
		return true
	}
	for _, val := range c.allowedOrigins {
		if val == originHeader {
			return true
		}
	}
	return false
}

func (c *CorsHandler) isAllowedMethod(methodHeader string) bool {
	if len(c.allowedMethods) == 0 {
		return false
	}
	if methodHeader == "OPTIONS" {
		return true
	}
	for _, m := range c.allowedMethods {
		if m == methodHeader {
			return true
		}
	}
	return false
}

func (c *CorsHandler) areHeadersAllowed(headers []string) bool {
	if c.allowedHeadersAll || len(headers) == 0 {
		return true
	}
	for _, header := range headers {
		found := false
		for _, h := range c.allowedHeaders {
			if h == header {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
