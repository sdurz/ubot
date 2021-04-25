package ubot

import (
	"reflect"
	"testing"

	"github.com/sdurz/axon"
)

func Test_And(t *testing.T) {
	type args struct {
		matchers []UMatcher
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "and ok",
			args: args{
				matchers: []UMatcher{
					func(b *Bot, o axon.O) bool {
						return true
					},
					func(b *Bot, o axon.O) bool {
						return true
					},
				},
			},
			want: true,
		},
		{
			name: "and ko",
			args: args{
				matchers: []UMatcher{
					func(b *Bot, o axon.O) bool {
						return true
					},
					func(b *Bot, o axon.O) bool {
						return false
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := And(tt.args.matchers...)(nil, axon.O{}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("And() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Or(t *testing.T) {
	type args struct {
		matchers []UMatcher
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "or ok",
			args: args{
				matchers: []UMatcher{
					func(b *Bot, o axon.O) bool {
						return false
					},
					func(b *Bot, o axon.O) bool {
						return true
					},
				},
			},
			want: true,
		},
		{
			name: "or ko",
			args: args{
				matchers: []UMatcher{
					func(b *Bot, o axon.O) bool {
						return false
					},
					func(b *Bot, o axon.O) bool {
						return false
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Or(tt.args.matchers...)(nil, axon.O{}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("And() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsFrom(t *testing.T) {
	type args struct {
		userID int64
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		message map[string]interface{}
	}{
		{
			name: "private match",
			args: args{1234},
			want: true,
			message: map[string]interface{}{
				"from": map[string]interface{}{
					"id": 1234.,
				},
			},
		},
		{
			name: "private nomatch",
			args: args{12345},
			want: false,
			message: map[string]interface{}{
				"from": map[string]interface{}{
					"id": 1234.,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFrom(tt.args.userID)(nil, tt.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_MessageHasCommand(t *testing.T) {
	type args struct {
		bot    *Bot
		entity string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		message map[string]interface{}
	}{
		{
			name: "command match",
			args: args{nil, "/cmd"},
			want: true,
			message: map[string]interface{}{
				"chat": map[string]interface{}{
					"type": "private",
				},
				"text": "123/cmd",
				"entities": []interface{}{
					map[string]interface{}{
						"offset": 3.,
						"length": 4.,
					},
				},
			},
		},
		{
			name: "command match in group",
			args: args{
				bot: &Bot{
					BotUser: User{
						Username: "testuser",
					},
				},
				entity: "/cmd",
			},
			want: true,
			message: map[string]interface{}{
				"chat": map[string]interface{}{
					"type": "group",
				},
				"text": "123/cmd@testuser",
				"entities": []interface{}{
					map[string]interface{}{
						"offset": 3.,
						"length": 13.,
					},
				},
			},
		},
		{
			name: "command don't extsts",
			args: args{nil, "/cmd"},
			want: false,
			message: map[string]interface{}{
				"chat": map[string]interface{}{
					"type": "private",
				},
				"text": "123/entity",
				"entities": []interface{}{
					map[string]interface{}{
						"offset": 3.,
						"length": 7.,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MessageHasCommand(tt.args.entity)(tt.args.bot, tt.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MessageHasCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_MessageHasPhoto(t *testing.T) {
	type args struct {
		bot     *Bot
		message axon.O
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
	}{
		{
			name: "ok",
			args: args{
				nil,
				map[string]interface{}{
					"photo": []interface{}{1, 2, 3},
				},
			},
			wantResult: true,
		},
		{
			name: "ok",
			args: args{
				nil,
				map[string]interface{}{
					"photoz": []interface{}{1, 2, 3},
				},
			},
			wantResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := MessageHasPhoto(tt.args.bot, tt.args.message); gotResult != tt.wantResult {
				t.Errorf("MessageHasPhoto() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_MessageHasEntities(t *testing.T) {
	type args struct {
		bot     *Bot
		message axon.O
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
	}{
		{
			name: "ok",
			args: args{
				nil,
				map[string]interface{}{
					"entities": []interface{}{1, 2, 3},
				},
			},
			wantResult: true,
		},
		{
			name: "ok",
			args: args{
				nil,
				map[string]interface{}{
					"photo": []interface{}{1, 2, 3},
				},
			},
			wantResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := MessageHasEntities(tt.args.bot, tt.args.message); gotResult != tt.wantResult {
				t.Errorf("MessageHasEntities() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_MessageInGroup(t *testing.T) {
	type args struct {
		bot     *Bot
		message axon.O
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
	}{
		{
			name: "ok group",
			args: args{
				nil,
				map[string]interface{}{
					"photo": []interface{}{1, 2, 3},
					"chat": map[string]interface{}{
						"type": "group",
					},
				},
			},
			wantResult: true,
		},
		{
			name: "ok supergroup",
			args: args{
				nil,
				map[string]interface{}{
					"photo": []interface{}{1, 2, 3},
					"chat": map[string]interface{}{
						"type": "supergroup",
					},
				},
			},
			wantResult: true,
		},
		{
			name: "ko private",
			args: args{
				nil,
				map[string]interface{}{
					"photo": []interface{}{1, 2, 3},
					"chat": map[string]interface{}{
						"type": "private",
					},
				},
			},
			wantResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := MessageInGroup(tt.args.bot, tt.args.message); gotResult != tt.wantResult {
				t.Errorf("MessageInGroup() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
