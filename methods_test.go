package ubot

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/sdurz/axon"
)

type mockAPIClient struct {
	method          string
	bytesMethod     func() []byte
	interfaceMethod func() interface{}
}

func (m *mockAPIClient) GetBytes(URL string) (result []byte, err error) {
	if !strings.Contains(URL, m.method) {
		err = errors.New("path don't match method")
	} else {
		result = m.bytesMethod()
	}
	return
}

func (m *mockAPIClient) PostBytes(URL string, data interface{}) (result []byte, err error) {
	if !strings.Contains(URL, m.method) {
		err = errors.New("path don't match method")
	} else {
		result = m.bytesMethod()
	}
	return
}

func (m *mockAPIClient) GetJson(URL string) (result interface{}, err error) {
	if !strings.Contains(URL, m.method) {
		err = errors.New("path don't match method")
	} else {
		result = m.interfaceMethod()
	}
	return
}
func (m *mockAPIClient) PostJson(URL string, request interface{}) (result interface{}, err error) {
	if !strings.Contains(URL, m.method) {
		err = errors.New("path don't match method")
	} else {
		result = m.interfaceMethod()
	}
	return
}
func (m *mockAPIClient) PostMultipart(URL string, request axon.O) (result interface{}, err error) {
	if !strings.Contains(URL, m.method) {
		err = errors.New("path don't match method")
	} else {
		result = m.interfaceMethod()
	}
	return
}

func TestBot_GetMe(t *testing.T) {
	type fields struct {
		Configuration Configuration
		apiClient     apiClient
		BotUser       User
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult User
		wantErr    bool
	}{
		{
			name: "test1",
			fields: fields{
				apiClient: &mockAPIClient{
					method: "getMe",
					interfaceMethod: func() interface{} {
						return map[string]interface{}{
							"id":                          123456.,
							"is_bot":                      true,
							"first_name":                  "A",
							"last_name":                   "B",
							"username":                    "botuser",
							"language_code":               "it",
							"can_join_groups":             true,
							"can_read_all_group_messages": true,
							"supports_inline_queries":     true,
						}
					},
					bytesMethod: func() []byte {
						return []byte("")
					},
				},
			},
			wantResult: User{
				ID:                      123456,
				IsBot:                   true,
				FirstName:               "A",
				LastName:                "B",
				Username:                "botuser",
				LanguageCode:            "it",
				CanJoinGroups:           true,
				CanReadAllGroupMessages: true,
				SupportsInlineQueries:   true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				Configuration: tt.fields.Configuration,
				apiClient:     tt.fields.apiClient,
				BotUser:       tt.fields.BotUser,
			}
			gotResult, err := b.GetMe()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.GetMe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*gotResult, tt.wantResult) {
				t.Errorf("Bot.GetMe() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestBot_LogOut(t *testing.T) {
	type fields struct {
		Configuration         Configuration
		apiClient             apiClient
		BotUser               User
		messageMHs            []matcherHandler
		editedMessageMHs      []matcherHandler
		channelPostMHs        []matcherHandler
		editedChannelPostMHs  []matcherHandler
		inlineQueryMHs        []matcherHandler
		chosenInlineResultMHs []matcherHandler
		callbackQueryMHs      []matcherHandler
		myChatMemberMHs       []matcherHandler
		chatMemberMHs         []matcherHandler
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{
				apiClient: &mockAPIClient{
					method: "logOut",
					interfaceMethod: func() interface{} {
						return map[string]interface{}{
							"id": 123456.,
						}
					},
					bytesMethod: func() []byte {
						return []byte("")
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				Configuration:         tt.fields.Configuration,
				apiClient:             tt.fields.apiClient,
				BotUser:               tt.fields.BotUser,
				messageMHs:            tt.fields.messageMHs,
				editedMessageMHs:      tt.fields.editedMessageMHs,
				channelPostMHs:        tt.fields.channelPostMHs,
				editedChannelPostMHs:  tt.fields.editedChannelPostMHs,
				inlineQueryMHs:        tt.fields.inlineQueryMHs,
				chosenInlineResultMHs: tt.fields.chosenInlineResultMHs,
				callbackQueryMHs:      tt.fields.callbackQueryMHs,
				myChatMemberMHs:       tt.fields.myChatMemberMHs,
				chatMemberMHs:         tt.fields.chatMemberMHs,
			}
			if err := b.LogOut(); (err != nil) != tt.wantErr {
				t.Errorf("Bot.LogOut() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBot_Close(t *testing.T) {
	type fields struct {
		Configuration         Configuration
		apiClient             apiClient
		BotUser               User
		messageMHs            []matcherHandler
		editedMessageMHs      []matcherHandler
		channelPostMHs        []matcherHandler
		editedChannelPostMHs  []matcherHandler
		inlineQueryMHs        []matcherHandler
		chosenInlineResultMHs []matcherHandler
		callbackQueryMHs      []matcherHandler
		myChatMemberMHs       []matcherHandler
		chatMemberMHs         []matcherHandler
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{
				apiClient: &mockAPIClient{
					method: "close",
					interfaceMethod: func() interface{} {
						return map[string]interface{}{
							"id": 123456.,
						}
					},
					bytesMethod: func() []byte {
						return []byte("")
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				Configuration:         tt.fields.Configuration,
				apiClient:             tt.fields.apiClient,
				BotUser:               tt.fields.BotUser,
				messageMHs:            tt.fields.messageMHs,
				editedMessageMHs:      tt.fields.editedMessageMHs,
				channelPostMHs:        tt.fields.channelPostMHs,
				editedChannelPostMHs:  tt.fields.editedChannelPostMHs,
				inlineQueryMHs:        tt.fields.inlineQueryMHs,
				chosenInlineResultMHs: tt.fields.chosenInlineResultMHs,
				callbackQueryMHs:      tt.fields.callbackQueryMHs,
				myChatMemberMHs:       tt.fields.myChatMemberMHs,
				chatMemberMHs:         tt.fields.chatMemberMHs,
			}
			if err := b.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Bot.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBot_SendMessage(t *testing.T) {
	type fields struct {
		Configuration         Configuration
		apiClient             apiClient
		BotUser               User
		messageMHs            []matcherHandler
		editedMessageMHs      []matcherHandler
		channelPostMHs        []matcherHandler
		editedChannelPostMHs  []matcherHandler
		inlineQueryMHs        []matcherHandler
		chosenInlineResultMHs []matcherHandler
		callbackQueryMHs      []matcherHandler
		myChatMemberMHs       []matcherHandler
		chatMemberMHs         []matcherHandler
	}
	type args struct {
		request axon.O
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult axon.O
		wantErr    bool
	}{
		{
			fields: fields{
				apiClient: &mockAPIClient{
					method: "sendMessage",
					interfaceMethod: func() interface{} {
						return map[string]interface{}{
							"id": 123456.,
						}
					},
					bytesMethod: func() []byte {
						return []byte("")
					},
				},
			},
			wantResult: axon.O{
				"id": 123456.,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				Configuration:         tt.fields.Configuration,
				apiClient:             tt.fields.apiClient,
				BotUser:               tt.fields.BotUser,
				messageMHs:            tt.fields.messageMHs,
				editedMessageMHs:      tt.fields.editedMessageMHs,
				channelPostMHs:        tt.fields.channelPostMHs,
				editedChannelPostMHs:  tt.fields.editedChannelPostMHs,
				inlineQueryMHs:        tt.fields.inlineQueryMHs,
				chosenInlineResultMHs: tt.fields.chosenInlineResultMHs,
				callbackQueryMHs:      tt.fields.callbackQueryMHs,
				myChatMemberMHs:       tt.fields.myChatMemberMHs,
				chatMemberMHs:         tt.fields.chatMemberMHs,
			}
			gotResult, err := b.SendMessage(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.SendMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Bot.SendMessage() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestBot_ForwardMessage(t *testing.T) {
	type fields struct {
		Configuration         Configuration
		apiClient             apiClient
		BotUser               User
		messageMHs            []matcherHandler
		editedMessageMHs      []matcherHandler
		channelPostMHs        []matcherHandler
		editedChannelPostMHs  []matcherHandler
		inlineQueryMHs        []matcherHandler
		chosenInlineResultMHs []matcherHandler
		callbackQueryMHs      []matcherHandler
		myChatMemberMHs       []matcherHandler
		chatMemberMHs         []matcherHandler
	}
	type args struct {
		request axon.O
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult axon.O
		wantErr    bool
	}{
		{
			fields: fields{
				apiClient: &mockAPIClient{
					method: "forwardMessage",
					interfaceMethod: func() interface{} {
						return map[string]interface{}{
							"id": 123456.,
						}
					},
					bytesMethod: func() []byte {
						return []byte("")
					},
				},
			},
			wantResult: axon.O{
				"id": 123456.,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				Configuration:         tt.fields.Configuration,
				apiClient:             tt.fields.apiClient,
				BotUser:               tt.fields.BotUser,
				messageMHs:            tt.fields.messageMHs,
				editedMessageMHs:      tt.fields.editedMessageMHs,
				channelPostMHs:        tt.fields.channelPostMHs,
				editedChannelPostMHs:  tt.fields.editedChannelPostMHs,
				inlineQueryMHs:        tt.fields.inlineQueryMHs,
				chosenInlineResultMHs: tt.fields.chosenInlineResultMHs,
				callbackQueryMHs:      tt.fields.callbackQueryMHs,
				myChatMemberMHs:       tt.fields.myChatMemberMHs,
				chatMemberMHs:         tt.fields.chatMemberMHs,
			}
			gotResult, err := b.ForwardMessage(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.SendMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Bot.SendMessage() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestBot_CopyMessage(t *testing.T) {
	type fields struct {
		Configuration Configuration
		apiClient     apiClient
	}
	type args struct {
		request axon.O
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult axon.O
		wantErr    bool
	}{
		{
			fields: fields{
				apiClient: &mockAPIClient{
					method: "copyMessage",
					interfaceMethod: func() interface{} {
						return map[string]interface{}{
							"id": 123456.,
						}
					},
					bytesMethod: func() []byte {
						return []byte("")
					},
				},
			},
			wantResult: axon.O{
				"id": 123456.,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				Configuration: tt.fields.Configuration,
				apiClient:     tt.fields.apiClient,
			}
			gotResult, err := b.CopyMessage(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.CopyMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Bot.CopyMessage() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestBot_GetMyCommands(t *testing.T) {
	type fields struct {
		Configuration Configuration
		apiClient     apiClient
		BotUser       User
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult axon.A
		wantErr    bool
	}{
		{
			name: "test1",
			fields: fields{
				apiClient: &mockAPIClient{
					method: "getMyCommands",
					interfaceMethod: func() interface{} {
						return []interface{}{1, 2, 3}
					},
					bytesMethod: func() []byte {
						return []byte("[1, 2, 3]")
					},
				},
			},
			wantResult: axon.A{
				1,
				2,
				3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				Configuration: tt.fields.Configuration,
				apiClient:     tt.fields.apiClient,
				BotUser:       tt.fields.BotUser,
			}
			gotResult, err := b.GetMyCommands()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.GetMyCommands() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Bot.GetMyCommands() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestBot_SetMyCommands(t *testing.T) {
	type fields struct {
		Configuration Configuration
		apiClient     apiClient
	}
	type args struct {
		request axon.O
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult bool
		wantErr    bool
	}{
		{
			name: "test1",
			fields: fields{
				apiClient: &mockAPIClient{
					method: "setMyCommands",
					interfaceMethod: func() interface{} {
						return true
					},
					bytesMethod: func() []byte {
						return []byte("true")
					},
				},
			},
			wantResult: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				Configuration: tt.fields.Configuration,
				apiClient:     tt.fields.apiClient,
			}
			gotResult, err := b.SetMyCommands(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.SetMyCommands() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("Bot.SetMyCommands() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestBot_AnswerCallbackQuery(t *testing.T) {
	type fields struct {
		apiClient apiClient
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult bool
		wantErr    bool
	}{
		{
			name: "test1",
			fields: fields{
				apiClient: &mockAPIClient{
					method: "answerCallbackQuery",
					interfaceMethod: func() interface{} {
						return true
					},
					bytesMethod: func() []byte {
						return []byte("true")
					},
				},
			},
			wantResult: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				apiClient: tt.fields.apiClient,
			}
			gotResult, err := b.AnswerCallbackQuery(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.AnswerCallbackQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("Bot.AnswerCallbackQuery() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestBot_DeleteChatStickerSet(t *testing.T) {
	type fields struct {
		apiClient apiClient
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult bool
		wantErr    bool
	}{
		{
			name: "test1",
			fields: fields{
				apiClient: &mockAPIClient{
					method: "deleteChatStickerSet",
					interfaceMethod: func() interface{} {
						return true
					},
					bytesMethod: func() []byte {
						return []byte("true")
					},
				},
			},
			wantResult: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				apiClient: tt.fields.apiClient,
			}
			gotResult, err := b.DeleteChatStickerSet(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.DeleteChatStickerSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("Bot.DeleteChatStickerSet() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestBot_SetChatStickerSet(t *testing.T) {
	type fields struct {
		apiClient apiClient
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult bool
		wantErr    bool
	}{
		{
			name: "test1",
			fields: fields{
				apiClient: &mockAPIClient{
					method: "setChatStickerSet",
					interfaceMethod: func() interface{} {
						return true
					},
					bytesMethod: func() []byte {
						return []byte("true")
					},
				},
			},
			wantResult: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				apiClient: tt.fields.apiClient,
			}
			gotResult, err := b.SetChatStickerSet(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.SetChatStickerSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("Bot.SetChatStickerSet() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestBot_GetChatMember(t *testing.T) {
	type fields struct {
		apiClient apiClient
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult axon.O
		wantErr    bool
	}{
		{
			name: "test1",
			fields: fields{
				apiClient: &mockAPIClient{
					method: "getChatMember",
					interfaceMethod: func() interface{} {
						return map[string]interface{}{
							"chat_id": 123456,
						}
					},
					bytesMethod: func() []byte {
						return []byte("true")
					},
				},
			},
			wantResult: map[string]interface{}{
				"chat_id": 123456,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				apiClient: tt.fields.apiClient,
			}
			gotResult, err := b.GetChatMember(axon.O{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.GetChatMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Bot.GetChatMember() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestBot_GetChatMembersCount(t *testing.T) {
	type fields struct {
		apiClient apiClient
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult int64
		wantErr    bool
	}{
		{
			name: "test1",
			fields: fields{
				apiClient: &mockAPIClient{
					method: "getChatMembersCount",
					interfaceMethod: func() interface{} {
						return 10.
					},
					bytesMethod: func() []byte {
						return []byte("10")
					},
				},
			},
			wantResult: 10,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				apiClient: tt.fields.apiClient,
			}
			gotResult, err := b.GetChatMembersCount(axon.O{"chat_id": 123456})
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.GetChatMembersCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Bot.GetChatMembersCount() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestBot_GetChatAdministrators(t *testing.T) {
	type fields struct {
		apiClient apiClient
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult axon.A
		wantErr    bool
	}{
		{
			name: "test1",
			fields: fields{
				apiClient: &mockAPIClient{
					method: "getChatAdministrators",
					interfaceMethod: func() interface{} {
						return []interface{}{
							map[string]interface{}{},
						}
					},
					bytesMethod: func() []byte {
						return []byte("[{}]")
					},
				},
			},
			wantResult: []interface{}{
				map[string]interface{}{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				apiClient: tt.fields.apiClient,
			}
			gotResult, err := b.GetChatAdministrators(axon.O{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.GetChatAdministrators() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Bot.GetChatAdministrators() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestBot_GetChat(t *testing.T) {
	type fields struct {
		apiClient apiClient
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult axon.O
		wantErr    bool
	}{
		{
			name: "test1",
			fields: fields{
				apiClient: &mockAPIClient{
					method: "getChat",
					interfaceMethod: func() interface{} {
						return make(map[string]interface{})
					},
					bytesMethod: func() []byte {
						return []byte("{}")
					},
				},
			},
			wantResult: map[string]interface{}{},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				apiClient: tt.fields.apiClient,
			}
			gotResult, err := b.GetChat(axon.O{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.GetChat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Bot.GetChat() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestBot_LeaveChat(t *testing.T) {
	type fields struct {
		apiClient apiClient
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult bool
		wantErr    bool
	}{
		{
			name: "test1",
			fields: fields{
				apiClient: &mockAPIClient{
					method: "leaveChat",
					interfaceMethod: func() interface{} {
						return true
					},
					bytesMethod: func() []byte {
						return []byte("true")
					},
				},
			},
			wantResult: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				apiClient: tt.fields.apiClient,
			}
			gotResult, err := b.LeaveChat(axon.O{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.LeaveChat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Bot.LeaveChat() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
