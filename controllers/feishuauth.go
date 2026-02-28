package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"webserver/middlewares"
	"webserver/models"
	"webserver/utils"

	"github.com/gin-gonic/gin"
)

// feishuUser 定义飞书用户信息结构体
type feishuUser struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`

	EnName          string `json:"en_name"`
	AvatarURL       string `json:"avatar_url"`
	AvatarThumb     string `json:"avatar_thumb"`
	AvatarMiddle    string `json:"avatar_middle"`
	AvatarBig       string `json:"avatar_big"`
	OpenID          string `json:"open_id"`
	UnionID         string `json:"union_id"`
	Email           string `json:"email"`
	EnterpriseEmail string `json:"enterprise_email"`
	TenantKey       string `json:"tenant_key"`
	EmployeeNo      string `json:"employee_no"`
}

// FeishuTokenResponse 定义获取 token 的响应结构
type FeishuTokenResponse struct {
	Code         int    `json:"code"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// FeishuUserInfoResponse 定义获取用户信息的响应结构
type FeishuUserInfoResponse struct {
	Code int         `json:"code"`
	Data *feishuUser `json:"data"`
	Msg  string      `json:"msg"`
}

type FeishuAuth struct{}

func (x *FeishuAuth) URLPatterns() []utils.Route {
	return []utils.Route{
		{Method: http.MethodGet, Path: "/flogin", ResourceFunc: x.Login},
	}
}

func (it *FeishuAuth) Login(c *gin.Context) {
	code := c.Query("code")
	errorStr := c.Query("error")

	if errorStr != "" {
		utils.ErrorMsg(c, errorStr)
		return
	}

	if code == "" {
		utils.ErrorMsg(c, "no code")
		return
	}

	// 请求地址
	url := "https://open.feishu.cn/open-apis/authen/v2/oauth/token"

	// 请求数据
	data := map[string]interface{}{
		"grant_type":    "authorization_code",
		"code":          code,
		"client_id":     utils.AppConfig.Feishu.ClientID,
		"client_secret": utils.AppConfig.Feishu.ClientSecret,
		"redirect_uri":  utils.AppConfig.Feishu.RedirectURI,
	}

	jsonData, _ := json.Marshal(data)

	// 创建请求
	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		utils.ErrorMsg(c, fmt.Sprintf("failed to create request: %v", err))
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.ErrorMsg(c, fmt.Sprintf("Curl error: %v", err))
		return
	}
	defer resp.Body.Close()

	var tokenResp FeishuTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		utils.ErrorMsg(c, fmt.Sprintf("failed to decode response: %v", err))
		return
	}

	if tokenResp.Code != 0 {
		utils.ErrorMsg(c, fmt.Sprintf("token response error: %s", resp.Body))
		return
	}

	// 获取用户信息
	userInfoURL := "https://open.feishu.cn/open-apis/authen/v1/user_info"
	userInfoReq, err := http.NewRequestWithContext(context.Background(), "GET", userInfoURL, nil)
	if err != nil {
		utils.ErrorMsg(c, fmt.Sprintf("failed to create user info request: %v", err))
		return
	}

	userInfoReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenResp.AccessToken))
	userInfoReq.Header.Set("Content-Type", "application/json; charset=utf-8")

	userInfoResp, err := client.Do(userInfoReq)
	if err != nil {
		utils.ErrorMsg(c, fmt.Sprintf("Curl error: %v", err))
		return
	}
	defer userInfoResp.Body.Close()

	var userInfo FeishuUserInfoResponse
	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		utils.ErrorMsg(c, fmt.Sprintf("failed to decode user info: %v", err))
		return
	}

	if userInfo.Code != 0 {
		utils.ErrorMsg(c, fmt.Sprintf("user info error: %s", userInfo.Msg))
		return
	}

	// 检查用户并设置 session
	user, err := it.checkUser(userInfo.Data)
	if err != nil {
		utils.ErrorMsg(c, fmt.Sprintf("failed to check user: %v", err))
		return
	}

	middlewares.SaveSession(c, user.Id)

	c.Redirect(http.StatusMovedPermanently, "/task/ilist")
}

func (it *FeishuAuth) checkUser(userInfo *feishuUser) (*models.User, error) {
	user := &models.User{}
	has, err := models.DB.Where("account = ?", userInfo.UserID).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		user = &models.User{
			Account: userInfo.UserID,
			Name:    userInfo.Name,
		}

		user.Department = models.DepartmentList[0].Id

		// 生成随机密码
		user.Password = utils.RandNumeric(16)

		_, err := models.DB.Insert(user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}
