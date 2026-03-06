package models

import (
	"errors"
	"strconv"
	"strings"
	"webserver/utils"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id         int64  `json:"id" xorm:"pk autoincr notnull comment('主键')"`
	Account    string `json:"account" xorm:"notnull unique comment('账号')"`
	Name       string `json:"name" xorm:"notnull unique comment('姓名')"`
	Nick       string `json:"nick"`
	Password   string `json:"password" xorm:"notnull comment('密码')"`
	Department int64  `json:"department" xorm:"notnull comment('部门')"`
	Team       int    `json:"team" xorm:"notnull comment('用户组')"`
	IsAdmin    bool   `json:"is_admin" xorm:"tinyint(1) notnull default(0) comment('是否管理员')"`
	CreatedAt  int64  `json:"created_at" xorm:"notnull comment('创建时间') created"`
	UpdatedAt  int64  `json:"updated_at" xorm:"notnull comment('更新时间') updated"`
	Ps         string `json:"ps"`       //权限
	IsLeave    bool   `json:"is_leave"` //是否离职

	Token   string `json:"token" xorm:"-"`
	PsGroup []int  `json:"ps_group" xorm:"-"`
}

func (u *User) SaveUser() error {
	_, err := DB.InsertOne(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GenerateFromPassword() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
}

func (u *User) BeforeInsert() {
	u.BeforeUpdate()
}

func (u *User) BeforeUpdate() {
	u.Account = strings.TrimSpace(u.Account)
	u.Name = strings.TrimSpace(u.Name)
	u.Ps = ""
	for _, ps := range u.PsGroup {
		if u.Ps != "" {
			u.Ps += ","
		}
		u.Ps += strconv.Itoa(ps)
	}
}

func (u *User) AfterLoad() {
	ps := strings.Split(u.Ps, ",")
	u.PsGroup = make([]int, 0)
	if len(ps) == 0 {
		u.PsGroup = append(u.PsGroup, 0)
		return
	}
	for _, ps := range ps {
		psInt, _ := strconv.Atoi(ps)
		u.PsGroup = append(u.PsGroup, psInt)
	}
}

// 对哈希加密的密码进行比对校验
func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(account, password string) (*User, error) {
	u := &User{}

	has, err := DB.Where("account = ?", account).Get(u)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("user not found")
	}

	err = VerifyPassword(password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, err
	}

	token, err := utils.GenerateToken(u.Id)
	if err != nil {
		return nil, err
	}
	u.Token = token
	return u, nil
}

func GetUserByID(id int64) (*User, error) {
	u := new(User)
	has, err := DB.ID(id).Get(u)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("user not found")
	}

	return u, nil
}

func CheckUserExist(id int) bool {
	u := new(User)
	has, err := DB.Where("id = ?", id).Get(u)
	if err != nil {
		return false
	}
	if !has {
		return false
	}
	return true
}

func GetUsers() (map[int64]*User, error) {
	users := make(map[int64]*User, 0)
	err := DB.Find(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
