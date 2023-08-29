package entity

import "time"

type Operation struct {
	Id int `db:"id"`
	UserId int `db:"user_id"`
	SlugId int `db:"slug_id"`
	AddedAt time.Time `db:"added_at"`
	RemovedAt time.Time `db:"removed_at"`
}
