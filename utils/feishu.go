package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	uuid "github.com/google/uuid"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type SendMsg struct {
	ReceiveId string
	Content   string
}

var chMsgReceiver = make(chan *SendMsg, 10000)

// 1000 次/分钟、50 次/秒
// https://open.feishu.cn/document/server-docs/im-v1/message/create
func msgReceiver() {
	defer LogRecover()
	for msg := range chMsgReceiver {
		FeishuSender(msg)
		time.Sleep(25 * time.Millisecond)
	}
}

func init() {
	go msgReceiver()
}

// 将消息先存放到队列
func FeishuAddMsgToChannel(receiveId string, content string) {
	if receiveId == "" {
		return
	}

	sendMsg := &SendMsg{Content: content, ReceiveId: receiveId}
	chMsgReceiver <- sendMsg
	// FeishuSender(sendMsg)
}

func FeishuSender(sendMsg *SendMsg) {
	uuidStr := uuid.New().String()
	// 创建请求对象
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`user_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(sendMsg.ReceiveId).
			MsgType(`interactive`). //固定为卡片消息
			Content(sendMsg.Content).
			Uuid(uuidStr).
			Build()).
		Build()

	// 创建 Client
	client := lark.NewClient(AppConfig.Feishu.ClientID, AppConfig.Feishu.ClientSecret)

	// 发起请求
	resp, err := client.Im.V1.Message.Create(context.Background(), req)

	// 处理错误
	if err != nil {
		GetLogger().Errorf("logId: %s, error: %s", resp.RequestId(), err)
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		GetLogger().Errorf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return
	}

	// 业务处理
	// fmt.Println(larkcore.Prettify(resp))
}

type FeishuWebhookResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 发送飞书群
func FeishuWebhook(body string, id string) {
	resp, err := http.Post("https://open.feishu.cn/open-apis/bot/v2/hook/"+id,
		"application/json", strings.NewReader(body))
	if err != nil {
		GetLogger().Errorf("Post error: %v", err)
		return
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)

	data := FeishuWebhookResp{}
	json.Unmarshal(bodyBytes, &data)
	if data.Code != 0 {
		GetLogger().Errorf("Webhook error: %s", bodyBytes)
	}
}

const tmp_text = `
{
    "schema": "2.0",
    "config": {
        "update_multi": true,
        "style": {
            "text_size": {
                "normal_v2": {
                    "default": "normal",
                    "pc": "normal",
                    "mobile": "heading"
                }
            }
        }
    },
	"header": {
        "title": {
            "tag": "plain_text",
            "content": "%s"
        },
        "subtitle": {
            "tag": "plain_text",
            "content": ""
        },
        "template": "default",
        "padding": "12px 12px 12px 12px"
    },
    "body": {
        "direction": "vertical",
        "padding": "12px 12px 12px 12px",
        "elements": [
            {
                "tag": "div",
                "text": {
                    "tag": "plain_text",
                    "content": "%s",
                    "text_size": "normal_v2",
                    "text_align": "left",
                    "text_color": "default"
                },
                "margin": "0px 0px 0px 0px"
            }
        ]
    }
}
`

const webhook_text_tmp = `
{
  "msg_type": "interactive",
  "card": {
    "schema": "2.0",
    "config": {
        "update_multi": true,
        "style": {
            "text_size": {
                "normal_v2": {
                    "default": "normal",
                    "pc": "normal",
                    "mobile": "heading"
                }
            }
        }
    },
    "body": {
        "direction": "vertical",
        "padding": "12px 12px 12px 12px",
        "elements": [
            {
                "tag": "markdown",
                "content": "%s",
                "text_align": "left",
                "text_size": "normal_v2",
                "margin": "0px 0px 0px 0px"
            }
        ]
    }
  }
}
	`

func FeishuText(user, title, content string) {
	// 转义特殊字符，防止破坏JSON结构或Markdown渲染
	escapedContent := strings.ReplaceAll(content, "\"", "\\\"")
	escapedContent = strings.ReplaceAll(escapedContent, "\n", "\\n")
	escapedContent = strings.ReplaceAll(escapedContent, "\r", "\\r")

	s := fmt.Sprintf(tmp_text, title, escapedContent)
	FeishuAddMsgToChannel(user, s)
}

func FeishuTextWebhook(hookid, content string) {
	// 转义特殊字符，防止破坏JSON结构或Markdown渲染
	escapedContent := strings.ReplaceAll(content, "\"", "\\\"")
	escapedContent = strings.ReplaceAll(escapedContent, "\n", "\\n")
	escapedContent = strings.ReplaceAll(escapedContent, "\r", "\\r")

	s := fmt.Sprintf(webhook_text_tmp, escapedContent)
	FeishuWebhook(s, hookid)
}
