package ubot

import (
	"context"
	"reflect"
	"testing"

	"github.com/sdurz/axon"
)

func TestBot_process(t *testing.T) {
	type fields struct {
		Configuration         Configuration
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
		ctx    context.Context
		update axon.O
	}

	var invocations int
	mhStop := matcherHandler{
		matcher: func(*Bot, axon.O) bool {
			return true
		},
		handler: func(context.Context, *Bot, axon.O) (bool, error) {
			invocations++
			return true, nil
		},
	}
	mhContinue := matcherHandler{
		matcher: func(*Bot, axon.O) bool {
			return true
		},
		handler: func(context.Context, *Bot, axon.O) (bool, error) {
			invocations++
			return false, nil
		},
	}

	tests := []struct {
		name            string
		fields          fields
		args            args
		wantInvocations int
		wantErr         bool
	}{
		{
			name: "a base message test",
			fields: fields{
				messageMHs: []matcherHandler{mhStop},
			},
			args: args{
				update: axon.O{
					"message": map[string]interface{}{},
				},
			},
			wantInvocations: 1,
			wantErr:         false,
		},
		{
			name: "stop execution test",
			fields: fields{
				messageMHs: []matcherHandler{mhStop, mhContinue},
			},
			args: args{
				update: axon.O{
					"message": map[string]interface{}{},
				},
			},
			wantInvocations: 1,
			wantErr:         false,
		},
		{
			name: "continue execution test",
			fields: fields{
				messageMHs: []matcherHandler{mhContinue, mhContinue},
			},
			args: args{
				update: axon.O{
					"message": map[string]interface{}{},
				},
			},
			wantInvocations: 2,
			wantErr:         false,
		},
		{
			name: "a base edited_message test",
			fields: fields{
				editedMessageMHs: []matcherHandler{mhStop},
			},
			args: args{
				update: axon.O{
					"edited_message": map[string]interface{}{},
				},
			},
			wantInvocations: 1,
			wantErr:         false,
		},
		{
			name: "a base channel_post test",
			fields: fields{
				channelPostMHs: []matcherHandler{mhStop},
			},
			args: args{
				update: axon.O{
					"channel_post": map[string]interface{}{},
				},
			},
			wantInvocations: 1,
			wantErr:         false,
		},
		{
			name: "a base edited_channel_post test",
			fields: fields{
				editedChannelPostMHs: []matcherHandler{mhStop},
			},
			args: args{
				update: axon.O{
					"edited_channel_post": map[string]interface{}{},
				},
			},
			wantInvocations: 1,
			wantErr:         false,
		},
		{
			name: "a base callback_query test",
			fields: fields{
				callbackQueryMHs: []matcherHandler{mhStop},
			},
			args: args{
				update: axon.O{
					"callback_query": map[string]interface{}{},
				},
			},
			wantInvocations: 1,
			wantErr:         false,
		},
		{
			name: "a base inline_query test",
			fields: fields{
				inlineQueryMHs: []matcherHandler{mhStop},
			},
			args: args{
				update: axon.O{
					"inline_query": map[string]interface{}{},
				},
			},
			wantInvocations: 1,
			wantErr:         false,
		},
		{
			name: "a base choosen_inline_result test",
			fields: fields{
				chosenInlineResultMHs: []matcherHandler{mhStop},
			},
			args: args{
				update: axon.O{
					"chosen_inline_result": map[string]interface{}{},
				},
			},
			wantInvocations: 1,
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invocations = 0
			b := &Bot{
				Configuration:         tt.fields.Configuration,
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
			if err := b.process(tt.args.ctx, tt.args.update); (err != nil) != tt.wantErr {
				t.Errorf("Bot.process() error = %v, wantErr %v", err, tt.wantErr)
			}
			if invocations != tt.wantInvocations {
				t.Errorf("Bot.process() error = %v, called %v", invocations, true)
			}
		})
	}
}

func TestBot_methodURL(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("method didn't panic")
		}
	}()

	b := &Bot{}
	b.methodURL("")
}

func Test_matcherHandler_evaluate(t *testing.T) {
	type fields struct {
		matcher Matcher
		handler Handler
	}
	type args struct {
		ctx     context.Context
		bot     *Bot
		message axon.O
	}
	bot := &Bot{}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult bool
		wantErr    bool
	}{
		{
			name: "match and return true",
			fields: fields{
				matcher: func(*Bot, axon.O) bool {
					return true
				},
				handler: func(context.Context, *Bot, axon.O) (bool, error) {
					return true, nil
				},
			},
			args: args{
				ctx:     context.Background(),
				bot:     bot,
				message: nil,
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "don;t match and would return true",
			fields: fields{
				matcher: func(*Bot, axon.O) bool {
					return false
				},
				handler: func(context.Context, *Bot, axon.O) (bool, error) {
					return true, nil
				},
			},
			args: args{
				ctx:     context.Background(),
				bot:     bot,
				message: nil,
			},
			wantResult: false,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &matcherHandler{
				matcher: tt.fields.matcher,
				handler: tt.fields.handler,
			}
			gotResult, err := m.evaluate(tt.args.ctx, tt.args.bot, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("matcherHandler.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("matcherHandler.evaluate() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_NewBot(t *testing.T) {
	type args struct {
		configuration *Configuration
	}
	tests := []struct {
		name       string
		args       args
		wantResult *Bot
	}{
		{
			name: "default workers",
			args: args{
				configuration: &Configuration{
					APIToken: "!23",
					WorkerNo: 0,
				},
			},
			wantResult: &Bot{
				Configuration: Configuration{
					APIToken: "!23",
					WorkerNo: 5,
				},
				apiClient: &httpApiClient{},
			},
		},
		{
			name: "exact workers",
			args: args{
				configuration: &Configuration{
					APIToken: "!23",
					WorkerNo: 10,
				},
			},
			wantResult: &Bot{
				Configuration: Configuration{
					APIToken: "!23",
					WorkerNo: 10,
				},
				apiClient: &httpApiClient{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := NewBot(tt.args.configuration); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("NewBot() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
