package ubot

import (
	"reflect"
	"testing"
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
					func(b *Bot, o O) bool {
						return true
					},
					func(b *Bot, o O) bool {
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
					func(b *Bot, o O) bool {
						return true
					},
					func(b *Bot, o O) bool {
						return false
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := And(tt.args.matchers...)(nil, O{}); !reflect.DeepEqual(got, tt.want) {
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
					func(b *Bot, o O) bool {
						return false
					},
					func(b *Bot, o O) bool {
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
					func(b *Bot, o O) bool {
						return false
					},
					func(b *Bot, o O) bool {
						return false
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Or(tt.args.matchers...)(nil, O{}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("And() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFrom(t *testing.T) {
	type args struct {
		userID   int64
		chatType string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		message map[string]interface{}
	}{
		{
			name: "private match",
			args: args{1234, "private"},
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
			args: args{12345, "private"},
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
			args: args{1234, "private"},
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
			if got := IsFrom(tt.args.userID, tt.args.chatType)(nil, tt.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageHasCommand(t *testing.T) {
	type args struct {
		entity   string
		chatType string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		message map[string]interface{}
	}{
		{
			name:    "command don't match",
			args:    args{"/cmd", "private"},
			want:    false,
			message: map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MessageHasCommand(tt.args.entity, tt.args.chatType)(nil, O{}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MessageHasCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageHasPhoto(t *testing.T) {
	type args struct {
		bot     *Bot
		message O
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
					"photo": &O{},
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
