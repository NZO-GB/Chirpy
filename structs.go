package main

import(
	_ "github.com/lib/pq"
	uuid "github.com/google/uuid"
	"time"

)

type UserRequest struct {
	Email				string			`json:"email"`
	Password			string			`json:"password"`
	Expires_seconds		int				`json:"expires_in_seconds"`
}

type UserJSON struct {
	ID 					uuid.UUID		`json:"id"`
	Created_at  		time.Time		`json:"created_at"`
	Updated_at			time.Time		`json:"updated_at"`
	Email 				string			`json:"email"`
	Password			string			`json:"password"`
	Token				string			`json:"token"`
}

type ChirpJSON struct {
	ID 					uuid.UUID		`json:"id"`
	Created_at  		time.Time		`json:"created_at"`
	Updated_at			time.Time		`json:"updated_at"`
	Body 				string			`json:"body"`
	User_id				uuid.UUID		`json:"user_id"`
}	
