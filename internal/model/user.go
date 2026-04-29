package model

import "time"

type User struct {
	UserDataID      string    `json:"user_data_id" db:"USER_DATA_ID"`
	UserName        string    `json:"username" db:"USER_NAME"`
	Password        string    `json:"password" db:"PASSWORD"`
	Email           string    `json:"email" db:"EMAIL"`
	Region          string    `json:"region" db:"REGION"`
	PreferredGenres []string  `json:"preferred_genres" db:"PREFERRED_GENRES"`
	AvatarURL       string    `json:"avatar_url" db:"AVATAR_URL"`
	Bio             string    `json:"bio" db:"BIO"`
	UsrCrt          string    `json:"usr_crt" db:"USR_CRT"`
	UsrUpd          string    `json:"usr_upd" db:"USR_UPD"`
	DtmCrt          time.Time `json:"dtm_crt" db:"DTM_CRT"`
	DtmUpd          time.Time `json:"dtm_upd" db:"DTM_UPD"`
}
