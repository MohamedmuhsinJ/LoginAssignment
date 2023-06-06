package controllers

import db "github.com/mohamedmuhsinJ/loginAssignment/Db"

func SyncDb() {

	db.Db.AutoMigrate(
		&User{},
	)
}
