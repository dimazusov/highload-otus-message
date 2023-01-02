package user

const TableName = "users"

const Female = 0
const Male = 1

type User struct {
	ID       uint   `json:"id" db:"id" gorm:"primary_key"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Name     string `json:"name" db:"name"`
	Surname  string `json:"surname" db:"surname"`
	Age      uint   `json:"age" db:"age"`
	Sex      bool   `json:"sex" db:"sex"`
	City     string `json:"city" db:"city" form:"city"`
	Interest string `json:"interest" db:"interests"`
}

func (m User) TableName() string {
	return TableName
}
