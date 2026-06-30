package main

import(
	_ "github.com/lib/pq"
	uuid "github.com/google/uuid"
	db "Chirpy/internal/database"
	"time"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
	dbQueries		*db.Queries
	platform		string
	secret			string
	polkaKey		string
	}


type UserRequest struct {
	Email				string			`json:"email"`
	Password			string			`json:"password"`
}

type UserJSON struct {
	ID 					uuid.UUID		`json:"id"`
	Created_at  		time.Time		`json:"created_at"`
	Updated_at			time.Time		`json:"updated_at"`
	Email 				string			`json:"email"`
	Password			string			`json:"password"`
	Token				string			`json:"token"`
	Refresh_Token		string			`json:"refresh_token"`
	Is_Chirpy_Red		bool			`json:"is_chirpy_red"`
}

type ChirpJSON struct {
	ID 					uuid.UUID		`json:"id"`
	Created_at  		time.Time		`json:"created_at"`
	Updated_at			time.Time		`json:"updated_at"`
	Body 				string			`json:"body"`
	User_id				uuid.UUID		`json:"user_id"`
}

type DataJSON struct {
	User_ID				uuid.UUID		`json:"user_id"`
}

type ChirpyRedJSON struct {
	Event				string			`json:"event"`
	Data				DataJSON 		`json:"data"`
}