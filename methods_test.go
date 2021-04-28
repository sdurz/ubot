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
		wantResult []interface{}
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
			wantResult: []interface{}{
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
