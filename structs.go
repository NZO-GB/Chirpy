package main

import(
	_ "github.com/lib/pq"
	uuid "github.com/google/uuid"
	"time"

)

type UserJSON struct {
	ID 			uuid.UUID	`json:"id"`
	Created_at  time.Time	`json:"created_at"`
	Updated_at	time.Time	`json:"updated_at"`
	Email 		string		`json:"email"`
}

type ChirpJSON struct {
	ID 			uuid.UUID	`json:"id"`
	Created_at  time.Time	`json:"created_at"`
	Updated_at	time.Time	`json:"updated_at"`
	Body 		string		`json:"body"`
	User_id		string		`json:"user_id"`
}
