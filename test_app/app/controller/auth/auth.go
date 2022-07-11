package auth

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"test/test_app/app/constants"
	"test/test_app/app/controller"
	models "test/test_app/app/db/dto"
	"test/test_app/app/model/request"
	response "test/test_app/app/model/response"
	"test/test_app/app/service/aws"
	"test/test_app/app/service/correlation"
	"test/test_app/app/service/logger"
	"test/test_app/app/service/supabase"
	"time"
)

type IAuthController interface {
	SignUp() gin.HandlerFunc
	Login() gin.HandlerFunc
}

type AuthController struct {
	AWSService      aws.IAwsService
	SupabaseService supabase.ISupabaseService
}

func NewAuthController(d aws.IAwsService, s supabase.ISupabaseService) IAuthController {
	return &AuthController{
		AWSService:      d,
		SupabaseService: s,
	}
}

// This is used to create deployment request
func (u AuthController) SignUp() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := correlation.WithReqContext(c)
		log := logger.Logger(ctx)
		dataFromBody := request.AuthLoginRequest{}
		err := json.NewDecoder(c.Request.Body).Decode(&dataFromBody)
		if err != nil {
			log.Errorf("Invalid Request Body", err)
			controller.RespondWithError(c, http.StatusBadRequest, "Invalid Request Body")
			return
		}
		pass, err := bcrypt.GenerateFromPassword([]byte(dataFromBody.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Errorf("Invalid Password Encryption", err)
			controller.RespondWithError(c, http.StatusBadRequest, "Invalid Password Encryption")
			return
		}

		dataFromBody.Password = string(pass)

		_, err = u.SupabaseService.AddNewUsers(ctx, constants.UserTable, dataFromBody.Email, dataFromBody.Password)
		if err != nil {
			log.Errorf("Error at Insert user", err)
			controller.RespondWithError(c, http.StatusBadRequest, "Error at Insert user")
			return
		}
		controller.ResponseMethod(c, http.StatusAccepted, response.ResponseV2{Success: true, Message: "Created Successfully", Data: ""})
	}
	return fn
}

func (u AuthController) Login() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := correlation.WithReqContext(c)
		log := logger.Logger(ctx)
		dataFromBody := request.AuthLoginRequest{}
		err := json.NewDecoder(c.Request.Body).Decode(&dataFromBody)
		if err != nil {
			log.Errorf("Invalid Request Body", err)
			controller.RespondWithError(c, http.StatusBadRequest, "Invalid Request Body")
			return
		}
		q := constants.UserTable + "?select=*&email=eq." + dataFromBody.Email

		userData, err := u.SupabaseService.GetUser(ctx, q)
		if err != nil {
			log.Errorf("Error at Db query while fetching user data", err)
			controller.RespondWithError(c, http.StatusBadRequest, "Error at Db query while fetching user data")
			return
		}
		errf := bcrypt.CompareHashAndPassword([]byte(userData[0].Password), []byte(dataFromBody.Password))
		if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword {
			log.Errorf("Invalid Password Encryption", err)
			controller.RespondWithError(c, http.StatusUnauthorized, "Invalid login credentials. Please try again")
			return
		}
		expiresAt := time.Now().Add(time.Minute * 100000).Unix()

		tk := &models.Token{
			UserId: userData[0].ID,
			Email:  userData[0].Email,
			StandardClaims: &jwt.StandardClaims{
				ExpiresAt: expiresAt,
			},
		}
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
		tokenString, _ := token.SignedString([]byte("secret"))
		var userToken = make(map[string]string)
		userToken["jwt"] = tokenString
		controller.ResponseMethod(c, http.StatusAccepted, response.ResponseV2{Success: true, Message: "Login Successfully", Data: userToken})
	}
	return fn
}
