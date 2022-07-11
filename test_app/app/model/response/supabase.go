package response

import "time"

type Userdata []struct {
	ID        string       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}


type Filedata []struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    string       `json:"user_id"`
	URL       string    `json:"url"`
	FileName  string    `json:"file_name"`
}