package ubot

import (
	"reflect"
	"testing"

	"github.com/sdurz/axon"
)

func TestAnd(t *testing.T) {
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

func TestOr(t *testing.T) {
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

func TestIsFrom(t *testing.T) {
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
				"chat": map[string]interface{}{
					"type": "private",
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
				"chat": map[string]interface{}{
					"type": "private",
				},
			},
		},
		{
			name: "public nomatch",
			args: args{1234},
			want: false,
			message: map[string]interface{}{
				"from": map[string]interface{}{
					"id": 1234.,
				},
				"chat": map[string]interface{}{
					"type": "public",
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

func TestMessageHasCommand(t *testing.T) {
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
					BotUser: UUser{
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

func TestMessageHasPhoto(t *testing.T) {
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
			name: "has photo",
			args: args{
				bot: &Bot{},
				message: map[string]interface{}{
					"photo": &axon.O{},
				},
			},
			wantResult: true,
		},
		{
			name: "has photo",
			args: args{
				bot:     &Bot{},
				message: map[string]interface{}{},
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
