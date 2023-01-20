package migration

import "github.com/tashfi04/printbin-server/models"

var Models []interface{}

func init() {
	Models = append(Models, models.User{})
	Models = append(Models, models.File{})
}
