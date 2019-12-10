package requests

/*
* 登录表单参数
 */
type LoginForm struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

/*
* 注册表单参数
 */
type RegisterForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password"  binding:"required,gt=6"`
	NickName string `form:"nickname" binding:"required"`
	Code     string `form:"code" binding:"required"`
	Captcha  string `form:"captcha" binding:"omitempty"`
	Key      string `form:"key" binding:"omitempty"`
	Ref      string `form:"ref"`
}
