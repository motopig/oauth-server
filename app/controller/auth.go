package controller

import (
	"github.com/motopig/oauth-server/app/common"
	"github.com/motopig/oauth-server/app/model"
	"github.com/motopig/oauth-server/app/requests"
	"github.com/motopig/oauth-server/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	"net/http"
)

func LoginHandler(c *gin.Context) {
	store, err := session.Start(nil, c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// 检查用户名密码是否正确
	var form requests.LoginForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusForbidden, &utils.ParamsErr)
		return
	}

	var user model.User
	isEmail := common.VerifyEmailFormat(form.UserName)
	var status bool
	if isEmail {
		status, err = model.WhereGet(&user, "email = ? ", form.UserName)
	} else {
		status, err = model.WhereGet(&user, "mobile = ? ", form.UserName)
	}

	if status == false {
		c.JSON(http.StatusForbidden, &utils.UserInvalid)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, &utils.Err)
		return
	}
	if pass := common.PasswordVerify(form.Password, user.Password); !pass {
		c.JSON(http.StatusForbidden, &utils.PasswordNotMatch)
		return
	}

	store.Set("LoggedInUserID", common.Int2str(user.Id))
	_ = store.Save()

	c.JSON(http.StatusOK, utils.Success)
}

func AuthHandler(c *gin.Context) {
	// 登录之后提示用户是否授权登录信息给第三方
	store, err := session.Start(nil, c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	// 检查登录是否成功 没有成功则跳转登录页面
	if _, ok := store.Get("LoggedInUserID"); !ok {
		c.JSON(http.StatusForbidden, utils.Err)
		return
	}
	c.JSON(http.StatusOK, utils.Success)
}

func RegisterHandler(c *gin.Context) {
	var form requests.RegisterForm

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusForbidden, &utils.ParamsErr)
		return
	}

	isEmail := common.VerifyEmailFormat(form.Username)
	if isEmail == false {
		// check captcha
		//if verifyCaptcha(form.Key, form.Captcha) == false {
		//	c.JSON(http.StatusOK, &utils.ErrCaptcha)
		//	return
		//}
	}

	var user model.User
	if isEmail {
		// check user exists
		user = model.User{
			UserBase: &model.UserBase{
				Email: form.Username,
			},
		}
	} else {
		user = model.User{
			UserBase: &model.UserBase{
				Mobile: form.Username,
			},
		}
	}

	counts, err := user.Count()

	if err != nil {
		c.JSON(http.StatusBadGateway, &utils.Err)
		return
	}

	if counts > 0 {
		c.JSON(http.StatusForbidden, &utils.UserExists)
		return
	}

	// check sms code
	//if verify(form.Code, form.Username) == false {
	//	c.JSON(http.StatusOK, &utils.VerifyCodeFail)
	//	return
	//}

	dbsession := model.DB().NewSession()
	defer dbsession.Close()

	err = dbsession.Begin()
	// 邀请参数
	var InviteUser model.User
	var has bool
	var er error
	if form.Ref != "" {
		has, er = dbsession.Where("invite_code = ?", form.Ref).Get(&InviteUser)
		if er != nil {

		}
		if has == true {
			user.Inviter = form.Ref
		}
	}

	if isEmail {
		user.Email = form.Username
		user.Mobile = ""
	} else {
		user.Email = ""
		user.Mobile = form.Username
	}
	hashPassword, err := common.PasswordHash(form.Password)
	if err != nil {
		c.JSON(http.StatusOK, &utils.RegisterFail)
		return
	}
	user.Password = hashPassword
	user.UserName = form.NickName
	user.RegisterIps = c.Request.RemoteAddr
	user.InviteCode = common.RandomStr(16)
	user.Langue = "ZH"
	_, err = dbsession.InsertOne(&user)
	if err != nil {
		_ = dbsession.Rollback()
		c.JSON(http.StatusOK, &utils.RegisterFail)
		return
	}

	err = dbsession.Commit()
	if err != nil {
		return
	}

	utils.RegisterSuccess.Data = gin.H{"id": user.Id}
	c.JSON(http.StatusOK, &utils.RegisterSuccess)
}
