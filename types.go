package ubot

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/sdurz/axon"
)

// Handler is a function that will handle an API update
type Handler func(context.Context, *Bot, axon.O) (bool, error)

// Matcher is a function that will decide wheter an update will be handled by a Matcher
type Matcher func(*Bot, axon.O) bool

// UpdatesSource are function that will get updates from the API server or any other source
// and publish them onto a channel.
// A proper UpdateSource will handle the context argument as needed.
type UpdatesSource func(*Bot, context.Context, chan axon.O)

// User struct stores user infos for the bot user
type User struct {
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

// UploadFile struct embeds binary data for sending binaries with send* methods.
type UploadFile interface {
	GetName() string
	GetReader() io.Reader
}

type readerUploadFile struct {
	fileName string
	reader   io.Reader
}

func (b *readerUploadFile) GetName() string {
	return b.fileName
}

func (b *readerUploadFile) GetReader() (result io.Reader) {
	return b.reader
}

// NewReaderUploadFile creates an UploadFile instance from a reader and a name
func NewReaderUploadFile(fileName string, reader io.Reader) (result UploadFile, err error) {
	if fileName == "" {
		err = errors.New("zero fileName")
		return
	}
	if reader == nil {
		err = errors.New("nil reader")
		return
	}
	result = &readerUploadFile{
		fileName: fileName,
		reader:   reader,
	}
	return
}

type bytesUploadFile struct {
	fileName  string
	fileBytes []byte
}

func (b *bytesUploadFile) GetName() string {
	return b.fileName
}

func (b *bytesUploadFile) GetReader() (result io.Reader) {
	return bytes.NewReader(b.fileBytes)
}

// NewBytesUploadFile creates an UploadFile instance from an array of bytes and a name
func NewBytesUploadFile(fileName string, fileBytes []byte) (result UploadFile, err error) {
	if fileName == "" {
		err = errors.New("zero fileName")
		return
	}
	if len(fileBytes) == 0 {
		err = errors.New("zero length or nil fileBytes")
		return
	}
	result = &bytesUploadFile{
		fileName:  fileName,
		fileBytes: fileBytes,
	}
	return
}
