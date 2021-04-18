package ubot

import (
	"log"
)

func And(matchers ...UMatcher) UMatcher {
	return func(b *Bot, u O) bool {
		for _, matcher := range matchers {
			if !matcher(b, u) {
				return false
			}
		}
		return true
	}
}

func Or(matchers ...UMatcher) UMatcher {
	return func(b *Bot, u O) bool {
		for _, matcher := range matchers {
			if matcher(b, u) {
				return true
			}
		}
		return false
	}
}

func Always(bot *Bot, update O) (result bool) {
	result = true
	return
}

// IsFrom matches private messages from a known userID
func IsFrom(userID int64, chatType string) UMatcher {
	return func(bot *Bot, message O) (result bool) {
		var (
			err        error
			iFrom      int64
			sType      string
			okChatType bool
		)
		if iFrom, err = message.GetInteger("from.id"); err != nil {
			return
		}
		if sType, err = message.GetString("chat.type"); err == nil {
			okChatType = sType == chatType
		}
		return iFrom == userID && okChatType
	}
}

// MessageHasPhoto matches if updates has a photo
func MessageHasPhoto(bot *Bot, message O) (result bool) {
	_, result = message["photo"]
	return
}

// MatchMessageEntities matches if update has message entities
func MessageHasEntities(bot *Bot, message O) (result bool) {
	if _, err := message.Get("entities"); err != nil {
		result = false
	}
	return
}

// MessageHasCommand matches if update has a certain message entity
func MessageHasCommand(entity string, chatType string) func(bot *Bot, message O) (result bool) {
	return func(b *Bot, message O) (result bool) {
		var (
			group     bool
			err       error
			mChatType string
			entities  A
			text      string
		)

		if mChatType, err = message.GetString("chat.type"); err != nil {
			return
		}
		if chatType != "" && chatType != mChatType {
			return
		}
		if entities, err = message.GetArray("entities"); err != nil {
			return
		}
		for _, ntt := range entities {
			var (
				ok     bool
				offset int64
				length int64
				oNtt   O
			)
			if oNtt, ok = ntt.(map[string]interface{}); !ok {
				log.Println("MessageHasCommand: entity not an O")
				return
			}
			if offset, err = oNtt.GetInteger("offset"); err != nil {
				return
			}
			if length, err = oNtt.GetInteger("length"); err != nil {
				return
			}
			text, _ = message.GetString("text")
			nttText := text[offset:length]
			if !group && nttText == entity {
				result = true
			} else if group && b.BotUser.Username != "" && nttText == entity+"@"+b.BotUser.Username {
				result = true
			}
		}
		return
	}
}
