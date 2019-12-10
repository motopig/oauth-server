package model

type User struct {
	*UserBase `xorm:"extends"`
	BaseModel `xorm:"extends"`
}

type UserBase struct {
	UserName    string `xorm:"username notnull char(100) default('') comment(用户名)" json:"username"`
	Mobile      string `xorm:"mobile index notnull char(100) default('') comment(登录手机号码)" json:"mobile"`
	Password    string `xorm:"password notnull char(100) default('') comment(登录密码)" json:"password"`
	Email       string `xorm:"email index notnull char(100) default('') comment(用户邮箱)" json:"email"`
	RegisterIps string `xorm:"register_ips comment(注册ip)" json:"register_ips"`
	Inviter     string `xorm:"inviter char(32) null default('') comment(邀请人的邀请码)" json:"inviter"`
	InviteCode  string `xorm:"invite_code char(32) notnull default('') comment(邀请码)" json:"invite_code"`
	Langue      string `xorm:"langue char(10) null default 'ZH' ENUM('ZH','EN') comment(用户语言)" json:"langue"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) Count() (int64, error) {
	return DB().Count(u)
}
