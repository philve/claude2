package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	http2 "github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/tls-client"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oldweipro/claude-to-chatgpt/global"
	"github.com/oldweipro/claude-to-chatgpt/model"
	"io"
	"time"
)

var (
	jar     = tls_client.NewCookieJar()
	options = []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(3600),
		tls_client.WithClientProfile(tls_client.Safari_Ipad_15_6),
		tls_client.WithNotFollowRedirects(),
		// create cookieJar instance and pass it as argument
		tls_client.WithCookieJar(jar),
		// Disable SSL verification
		tls_client.WithInsecureSkipVerify(),
	}
	client, _ = tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
)

func RequestClaudeToResponse(c *gin.Context, params *model.ChatMessageRequest, stream bool) {
	appendMessageApi := global.ServerConfig.Claude.BaseUrl + "/api/append_message"
	err := client.SetProxy(global.HttpProxy)
	if err != nil {
		return
	}
	// 设置两个参数
	newStringUuid := uuid.NewString()
	// TODO 判断是否出错
	_, err = CreateChatConversations(newStringUuid)
	if err != nil {
		return
	}
	params.ConversationUuid = newStringUuid
	params.OrganizationUuid = global.ServerConfig.Claude.OrganizationUuid
	// 发起请求
	marshal, err := json.Marshal(params)
	if err != nil {
		fmt.Println("Marshal err:", err)
	}
	request, err := http2.NewRequest(http2.MethodPost, appendMessageApi, bytes.NewBuffer(marshal))
	if err != nil {
		fmt.Println(err)
		return
	}
	SetHeaders(request)
	response, err := client.Do(request)
	reader := bufio.NewReader(response.Body)
	var originalResponse model.ChatMessageResponse
	var isRole = true
	if stream {
		// Response content type is text/event-stream
		c.Header("Content-Type", "text/event-stream")
	} else {
		// Response content type is application/json
		c.Header("Content-Type", "application/json")
	}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		if len(line) < 6 {
			continue
		}
		line = line[6:]
		if originalResponse.StopReason != "stop_sequence" {
			err = json.Unmarshal([]byte(line), &originalResponse)
			if err != nil {
				continue
			}
			completionResponse := model.ChatCompletionStreamResponse{
				ID:      "chatcmpl-7f1DmyzTWtiysnyfSS4i187kus2Ao",
				Object:  "chat.completion.chunk",
				Created: time.Now().Unix(),
				Model:   "gpt-3.5-turbo-0613",
				Choices: []model.ChatCompletionStreamChoice{
					{
						Index: 0,
						Delta: model.ChatCompletionStreamChoiceDelta{
							Content: originalResponse.Completion,
						},
						FinishReason: nil,
					},
				},
			}
			if originalResponse.Completion == "" {
				completionResponse.Choices[0].FinishReason = "stop"
			}
			if isRole {
				completionResponse.Choices[0].Delta.Role = "assistant"
			}
			resp, _ := json.Marshal(completionResponse)
			responseString := "data: " + string(resp) + "\n\n"
			_, err = c.Writer.WriteString(responseString)
			if err != nil {
				return
			}
			isRole = false
			c.Writer.Flush()
		} else {
			if stream {
				c.Writer.WriteString("data: [DONE]\n\n")
			}
		}

	}
	err = DeleteChatConversations(newStringUuid)
	if err != nil {
		fmt.Println("delete err:", newStringUuid)
	}
}

func CreateChatConversations(newStringUuid string) (model.ChatConversationResponse, error) {
	var chatConversationResponse model.ChatConversationResponse
	chatConversationsApi := global.ServerConfig.Claude.BaseUrl + "/api/organizations/" + global.ServerConfig.Claude.OrganizationUuid + "/chat_conversations"
	err := client.SetProxy(global.HttpProxy)
	if err != nil {
		return chatConversationResponse, err
	}
	conversation := model.NewChatConversationRequest(newStringUuid, "")
	marshal, err := json.Marshal(conversation)
	if err != nil {
		fmt.Println("Marshal err:", err)
	}
	request, err := http2.NewRequest(http2.MethodPost, chatConversationsApi, bytes.NewBuffer(marshal))

	if err != nil {
		return chatConversationResponse, err
	}
	SetHeaders(request)

	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return chatConversationResponse, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return chatConversationResponse, err
	}
	err = json.Unmarshal(body, &chatConversationResponse)
	if err != nil {
		return chatConversationResponse, err
	}
	return chatConversationResponse, err
}

func DeleteChatConversations(newStringUuid string) error {
	err := client.SetProxy(global.HttpProxy)
	if err != nil {
		return err
	}
	chatConversationsApi := global.ServerConfig.Claude.BaseUrl + "/api/organizations/" + global.ServerConfig.Claude.OrganizationUuid + "/chat_conversations/"
	request, err := http2.NewRequest(http2.MethodDelete, chatConversationsApi+newStringUuid, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	SetHeaders(request)

	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()
	//body, err := io.ReadAll(res.Body)
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	//fmt.Println("删除对话:", newStringUuid, string(body))
	if res.StatusCode != 200 {
		return errors.New("delete chat conversations err")
	}
	return nil
}

func GetOrganizations() ([]model.OrganizationsResponse, error) {
	err := client.SetProxy(global.HttpProxy)
	if err != nil {
		return nil, err
	}
	organizationsApi := global.ServerConfig.Claude.BaseUrl + "/api/organizations"
	request, err := http2.NewRequest(http2.MethodGet, organizationsApi, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	SetHeaders(request)
	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var response []model.OrganizationsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Unmarshal err:", err)
		return nil, err
	}
	return response, err
}

func SetHeaders(r *http2.Request) {
	r.Header.Add("Cookie", global.ServerConfig.Claude.SessionKey)
	r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "*/*")
	r.Header.Add("Host", "claude.ai")
	r.Header.Add("Connection", "keep-alive")
}
