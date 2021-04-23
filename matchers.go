package ubot

import (
	"log"

	"github.com/sdurz/axon"
)

func And(matchers ...UMatcher) UMatcher {
	return func(b *Bot, u axon.O) bool {
		for _, matcher := range matchers {
			if !matcher(b, u) {
				return false
			}
		}
		return true
	}
}

func Or(matchers ...UMatcher) UMatcher {
	return func(b *Bot, u axon.O) bool {
		for _, matcher := range matchers {
			if matcher(b, u) {
				return true
			}
		}
		return false
	}
}

func Always(bot *Bot, update axon.O) (result bool) {
	result = true
	return
}

// ChatType matches axon.A certain chat type
func ChatType(chatType string) UMatcher {
	return func(bot *Bot, message axon.O) (result bool) {
		var (
			err   error
			sType string
		)
		if sType, err = message.GetString("chat.type"); err == nil {
			result = sType == chatType
		}
		return
	}
}

// IsFrom matches private messages from axon.A known userID
func IsFrom(userID int64) UMatcher {
	return func(bot *Bot, message axon.O) (result bool) {
		var (
			err   error
			iFrom int64
		)
		if iFrom, err = message.GetInteger("from.id"); err == nil {
			result = iFrom == userID
		}
		return
	}
}

// MessageHasPhoto matches if updates has axon.A photo
func MessageHasPhoto(bot *Bot, message axon.O) (result bool) {
	_, result = message["photo"]
	return
}

// MatchMessageEntities matches if update has message entities
func MessageHasEntities(bot *Bot, message axon.O) (result bool) {
	if _, err := message.Get("entities"); err != nil {
		result = false
	}
	return
}

// MessageHasCommand matches if update has axon.A certain message entity
func MessageHasCommand(entity string) func(bot *Bot, message axon.O) (result bool) {
	return func(b *Bot, message axon.O) (result bool) {
		var (
			chatType string
			isGroup  bool
			err      error
			entities axon.A
			text     string
		)

		if entities, err = message.GetArray("entities"); err != nil {
			return
		}
		if chatType, err = message.GetString("chat.type"); err != nil {
			return
		}

		switch chatType {
		case "group":
			isGroup = true
		case "supergroup":
			isGroup = true
		default:
			isGroup = false
		}
		for _, ntt := range entities {
			var (
				ok     bool
				offset int64
				length int64
				oNtt   axon.O
			)
			if oNtt, ok = ntt.(map[string]interface{}); !ok {
				log.Println("MessageHasCommand: entity not an axon.O")
				return
			}
			if offset, err = oNtt.GetInteger("offset"); err != nil {
				return
			}
			if length, err = oNtt.GetInteger("length"); err != nil {
				return
			}
			if text, err = message.GetString("text"); err != nil {
				return
			}
			nttText := text[offset : offset+length]
			if !isGroup && nttText == entity {
				result = true
			} else if isGroup && b.BotUser.Username != "" && nttText == entity+"@"+b.BotUser.Username {
				result = true
			}
		}
		return
	}
}
