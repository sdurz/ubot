package ubot

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"

	"github.com/sdurz/axon"
)

type UHandler func(context.Context, *Bot, axon.O) (bool, error)
type UMatcher func(*Bot, axon.O) bool

type UpdatesSource func(*Bot, context.Context, chan axon.O)
type UpdateCallbackFunc func(axon.O) error

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
	CanJoinGroups           bool   `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries,omitempty"`
}

type UReply struct {
	Ok          bool            `json:"ok"`
	ErrorCode   int64           `json:"error_code,omitempty"`
	Description string          `json:"description,omitempty"`
	Result      json.RawMessage `json:"result"`
}
