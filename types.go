package ubot

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
)

type UHandler func(context.Context, *Bot, O) (bool, error)
type UMatcher func(*Bot, O) bool

type UpdateSourceFunc func(chan O, chan int, chan error)
type UpdateCallbackFunc func(O) error

type UMultipart interface {
	asMIMEPart(string, *multipart.Writer) (fw io.Writer, fr io.Reader, err error)
}

type UUser struct {
	ID                      int64  `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name,omitempty"`
	Username                string `json:"username,omitempty"`
	LanguageCode            string `json:"language_code,omitempty"`
	CanJoinGroup            string `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages string `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries   string `json:"supports_inline_queries,omitempty"`
}

type UReply struct {
	Ok          bool            `json:"ok"`
	ErrorCode   int64           `json:"error_code,omitempty"`
	Description string          `json:"description,omitempty"`
	Result      json.RawMessage `json:"result"`
}
