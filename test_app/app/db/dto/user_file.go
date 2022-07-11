package dbmodels

import "time"

const (
	TABLE_USER_FILE = "user_file"

	COLUMN_ID         = "id"
	COLUMN_USER_ID    = "user_id"
	COLUMN_FILE_NAME  = "file_name"
	COLUMN_URL        = "url"
	COLUMN_CREATED_AT = "created_at"
)

type UserFile struct {
	ID        *int       `gorm:"column:id" json:"id"`
	UserId    *int       `gorm:"column:user_id" json:"user_id"`
	FileName  *string    `gorm:"column:file_name" json:"file_name"`
	Url       *string    `gorm:"column:url" json:"url"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
}
