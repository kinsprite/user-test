package main

// UserInfo  from user test service
type UserInfo struct {
	ID    int    `json:"id"    db:"id"`
	Name  string `json:"name"  db:"name"`
	Email string `json:"email" db:"email"`
}
