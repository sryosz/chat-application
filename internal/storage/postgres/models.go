// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package postgres

type User struct {
	ID       int64
	Username string
	Email    string
	Password []byte
}
