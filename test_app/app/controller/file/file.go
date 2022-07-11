package file

import (
	"fmt"
	"net/http"
	"test/test_app/app/constants"
	"test/test_app/app/controller"
	dbmodels "test/test_app/app/db/dto"
	"test/test_app/app/model/response"
	"test/test_app/app/service/aws"
	"test/test_app/app/service/correlation"
	"test/test_app/app/service/logger"
	"test/test_app/app/service/supabase"

	"github.com/gin-gonic/gin"
)

type IFileController interface {
	UploadFile() gin.HandlerFunc
	GetFile() gin.HandlerFunc
}

type FileController struct {
	AWSService      aws.IAwsService
	SupabaseService supabase.ISupabaseService
}

func NewFileController(d aws.IAwsService, s supabase.ISupabaseService) IFileController {
	return &FileController{
		AWSService:      d,
		SupabaseService: s,
	}
}

func (u *FileController) UploadFile() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := correlation.WithReqContext(c)
		log := logger.Logger(ctx)

		claims, exist := c.Get(constants.CtxClaims)
		if !exist {
			log.Errorf("Auth Failed to parse Jwt")
			controller.RespondWithError(c, http.StatusUnauthorized, "Auth Failed to parse Jwt")
			return
		}
		auth := claims.(*dbmodels.Token)
		userId := auth.UserId

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 5*1024*1024)

		file, handler, err := c.Request.FormFile("file")
		if err != nil {
			log.Errorf("Error getting the file", err)
			controller.RespondWithError(c, http.StatusInternalServerError, "Error getting the file")
			return
		}

		defer file.Close()

		fmt.Printf("Uploaded file name: %+v\n", handler.Filename)
		fmt.Printf("Uploaded file size %+v\n", handler.Size)
		fmt.Printf("File mime type %+v\n", handler.Header)

		url, err := u.AWSService.UploadFile(ctx, handler.Filename, file)
		if err != nil {
			log.Errorf("Error uploading file", err)
			controller.RespondWithError(c, http.StatusInternalServerError, "Error uploading file")
			return
		}

		var filedataTransform = make(map[string]interface{})

		filedataTransform["user_id"] = userId
		filedataTransform["url"] = url
		filedataTransform["file_name"] = handler.Filename

		_, err = u.SupabaseService.AddNewFile(ctx, constants.FileTable, filedataTransform)
		if err != nil {
			log.Errorf("Error creating file", err)
			controller.RespondWithError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		controller.ResponseMethod(c, http.StatusAccepted, response.ResponseV2{Success: true, Message: "File Uploaded Successfully", Data: filedataTransform})
	}
	return fn
}

func (u *FileController) GetFile() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := correlation.WithReqContext(c)
		log := logger.Logger(ctx)

		claims, exist := c.Get(constants.CtxClaims)
		if !exist {
			log.Errorf("Auth Failed to parse Jwt")
			controller.RespondWithError(c, http.StatusUnauthorized, "Auth Failed to parse Jwt")
			return
		}
		auth := claims.(*dbmodels.Token)
		userId := auth.UserId
		q := constants.FileTable + "?select=*&user_id=eq." + userId
		userFiles, err := u.SupabaseService.GetFile(ctx, q)
		if err != nil {
			log.Errorf("Error creating file", err)
			controller.RespondWithError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		controller.ResponseMethod(c, http.StatusAccepted, response.ResponseV2{Success: true, Message: "User Files ", Data: userFiles})
	}
	return fn
}
