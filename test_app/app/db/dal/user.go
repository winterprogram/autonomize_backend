package dal

import (
	"context"
	"test/test_app/app/db"
	dbmodels "test/test_app/app/db/dto"
)

type IUserFileDal interface {
	CreateUserFile(ctx context.Context, userFiler dbmodels.UserFile) error
	GetUserFile(ctx context.Context, userID int) ([]dbmodels.UserFile, error)
}

type UserFileDal struct {
	DBService *db.DBService
}

func NewUserFileDal(dbService *db.DBService) IUserFileDal {
	return &UserFileDal{
		DBService: dbService,
	}
}

func (d *UserFileDal) CreateUserFile(ctx context.Context, userFiler dbmodels.UserFile) error {
	tx := d.DBService.GetDB()
	return tx.Table(dbmodels.TABLE_USER_FILE).Omit(dbmodels.COLUMN_ID, dbmodels.COLUMN_CREATED_AT).Create(&userFiler).Error
}

func (d *UserFileDal) GetUserFile(ctx context.Context, userID int) ([]dbmodels.UserFile, error) {
	tx := d.DBService.GetDB()
	var userFile []dbmodels.UserFile

	if err := tx.Table(dbmodels.TABLE_USER_FILE).Select("*").Where(dbmodels.COLUMN_USER_ID+" = ?", userID).Find(&userFile).Error; err != nil {
		return userFile, err
	}

	return userFile, nil
}
