package ubot

import (
	"errors"

	"github.com/sdurz/axon"
)

// GetMe returns basic information about the bot in form of a User object.
// see https://core.telegram.org/bots/api#getme
func (b *Bot) GetMe() (result *User, err error) {
	var (
		uResult User
		iResult interface{}
		oResult axon.O
		ok      bool
	)
	if iResult, err = b.doGet("getMe"); err == nil {
		if oResult, ok = iResult.(map[string]interface{}); !ok {
			err = errors.New("doGet returned unexpected type")
			return
		}
		uResult.ID, _ = oResult.GetInteger("id")
		uResult.IsBot, _ = oResult.GetBoolean("is_bot")
		uResult.FirstName, _ = oResult.GetString("first_name")
		uResult.LastName, _ = oResult.GetString("last_name")
		uResult.Username, _ = oResult.GetString("username")
		uResult.LanguageCode, _ = oResult.GetString("language_code")
		uResult.CanJoinGroups, _ = oResult.GetBoolean("can_join_groups")
		uResult.CanReadAllGroupMessages, _ = oResult.GetBoolean("can_read_all_group_messages")
		uResult.SupportsInlineQueries, _ = oResult.GetBoolean("supports_inline_queries")
	}
	return &uResult, err
}

// LogOut logs the bot out of the cloud Bot API server
// see https://core.telegram.org/bots/api#logout
func (b *Bot) LogOut() (err error) {
	_, err = b.doGet("logOut")
	return
}

// Close closea the bot instance
// see https://core.telegram.org/bots/api#close
func (b *Bot) Close() (err error) {
	_, err = b.doGet("close")
	return
}

// SendMessage sends a text message
// see https://core.telegram.org/bots/api#sendmessage
func (b *Bot) SendMessage(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("sendMessage", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// ForwardMessage forwards messages of any kind
// see https://core.telegram.org/bots/api#forwardmessage
func (b *Bot) ForwardMessage(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("forwardMessage", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// CopyMessage copy messages of any kind. The method is analogous to the method forwardMessage, but the copied message doesn't have a link to the original message.
// see https://core.telegram.org/bots/api#copymessage
func (b *Bot) CopyMessage(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("copyMessage", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// SendPhoto sends a photo
// see https://core.telegram.org/bots/api#sendphoto
func (b *Bot) SendPhoto(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendPhoto", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// SendAudio sends an audio
// see https://core.telegram.org/bots/api#sendaudio
func (b *Bot) SendAudio(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendAudio", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// SendVideo sends a video
// see https://core.telegram.org/bots/api#senddocument
func (b *Bot) SendDocument(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendDocument", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// SendVideo sends a video
// see https://core.telegram.org/bots/api#sendvideo
func (b *Bot) SendVideo(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendVideo", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// SendAnimation sends an animation
// see https://core.telegram.org/bots/api#sendanimation
func (b *Bot) SendAnimation(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendAnimation", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// SendVoice sends a voice
// see https://core.telegram.org/bots/api#sendvoice
func (b *Bot) SendVoice(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendVoice", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// SendVoice sends a video note
// see https://core.telegram.org/bots/api#sendvideonote
func (b *Bot) SendVideoNote(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendVideoNote", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// SendVoice sends a media group
// see https://core.telegram.org/bots/api#sendmediagroup
func (b *Bot) SendMediaGroup(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendMediaGroup", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// SendLocation sends a location
// see https://core.telegram.org/bots/api#sendlocation
func (b *Bot) SendLocation(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("sendLocation", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// EditMessageLiveLocation sends a location
// see https://core.telegram.org/bots/api#editmessagelivelocation
func (b *Bot) EditMessageLiveLocation(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("editMessageLiveLocation", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// StopMessageLiveLocation sends a location
// see https://core.telegram.org/bots/api#stopmessagelivelocation
func (b *Bot) StopMessageLiveLocation(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("stopMessageLiveLocation", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// SendVenue sends a venue
// see https://core.telegram.org/bots/api#sendvenue
func (b *Bot) SendVenue(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("sendVenue", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// SendContact sends a venue
// see https://core.telegram.org/bots/api#sendcontact
func (b *Bot) SendContact(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("sendContact", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// SendPoll sends a poll
// see https://core.telegram.org/bots/api#sendpoll
func (b *Bot) SendPoll(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("sendPoll", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// SendDice sends a dice
// see https://core.telegram.org/bots/api#senddice
func (b *Bot) SendDice(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("sendDice", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// SendChatAction sends a chat action
// see https://core.telegram.org/bots/api#sendchataction
func (b *Bot) SendChatAction(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("sendChatAction", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsBool()
	}
	return
}

// GetUserProfilePhotos gets user profile photos.
// see https://core.telegram.org/bots/api#getuserprofilephotos
func (b *Bot) GetUserProfilePhotos(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("getUserProfilePhotos", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// GetFile gets basic info about a file and prepare it for downloading.
// see https://core.telegram.org/bots/api#getfile
func (b *Bot) GetFile(fileId string) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doGet("getFile?file_id=" + fileId); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// KickChatMember kicks a user from a group, a supergroup or a channel.
// see https://core.telegram.org/bots/api#kickchatmember
func (b *Bot) KickChatMember(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("kickChatMember", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsBool()
	}
	return
}

// UnbanChatMember unban a previously kicked user in a supergroup or channel.
// see https://core.telegram.org/bots/api#unbanchatmember
func (b *Bot) UnbanChatMember(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("unbanChatMember", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsBool()
	}
	return
}

// RestrictChatMember unban a previously kicked user in a supergroup or channel.
// see https://core.telegram.org/bots/api#restrictchatmember
func (b *Bot) RestrictChatMember(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("restrictChatMember", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsBool()
	}
	return
}

// PromoteChatMember unban a previously kicked user in a supergroup or channel.
// see https://core.telegram.org/bots/api#promotechatmember
func (b *Bot) PromoteChatMember(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("promoteChatMember", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsBool()
	}
	return
}

// SetWebhook implements setWebhook from Telegram Bot API.
// see https://core.telegram.org/bots/api#setwebhook
func (b *Bot) SetWebhook(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("setWebhook", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsBool()
	}
	return
}

// DeleteWebhook implements deleteWebhook from Telegram Bot API
// see https://core.telegram.org/bots/api#deletewebhook
func (b *Bot) DeleteWebhook(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("deleteWebhook", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsBool()
	}
	return
}

// GetWebhookInfo get current webhook status
// https://core.telegram.org/bots/api#getwebhookinfo
func (b *Bot) GetWebhookInfo() (result axon.O, err error) {
	var response interface{}
	if response, err = b.doGet("getWebhookInfo"); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// PinChatMessage pins a message for the given chat
// https://core.telegram.org/bots/api#pinchatmessage
func (b *Bot) PinChatMessage(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("pinChatMessage", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsBool()
	}
	return
}

// UnpinChatMessage removes a message from the list of pinned messages in a chat.
// see https://core.telegram.org/bots/api#unpinchatmessage
func (b *Bot) UnpinChatMessage(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("unpinChatMessage", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsBool()
	}
	return
}

// UnpinAllChatMessages clears the list of pinned messages in a chat.
// see https://core.telegram.org/bots/api#unpinallchatmessages
func (b *Bot) UnpinAllChatMessages(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("unpinAllChatMessages", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsBool()
	}
	return
}

// LeaveChat leave a group, supergroup or channel.
// see https://core.telegram.org/bots/api#leavechat
func (b *Bot) LeaveChat(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("leaveChat", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsBool()
	}
	return
}

// GetChat get up to date information about the chat.
// see https://core.telegram.org/bots/api#getchat
func (b *Bot) GetChat(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("getChat", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// GetChatAdministrators get the number of members in a chat.
// see https://core.telegram.org/bots/api#getchatmemberscount
func (b *Bot) GetChatAdministrators(request axon.O) (result axon.A, err error) {
	var response interface{}
	if response, err = b.doPost("getChatAdministrators", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsArray()
	}
	return
}

// GetChatMembersCount get the number of members in a chat.
// see https://core.telegram.org/bots/api#getchatmemberscount
func (b *Bot) GetChatMembersCount(request axon.O) (result int64, err error) {
	var response interface{}
	if response, err = b.doPost("getChatMembersCount", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsInteger()
	}
	return
}

// GetChatMember gets information about a member of a chat.
// see https://core.telegram.org/bots/api#getchatmember
func (b *Bot) GetChatMember(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("getChatMember", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsObject()
	}
	return
}

// SetChatStickerSet  set a new group sticker set for a supergroup.
// see https://core.telegram.org/bots/api#setchatstickerset
func (b *Bot) SetChatStickerSet(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("setChatStickerSet", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsBool()
	}
	return
}

// DeleteChatStickerSet  set a new group sticker set for a supergroup.
// see https://core.telegram.org/bots/api#deletechatstickerset
func (b *Bot) DeleteChatStickerSet(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("deleteChatStickerSet", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsBool()
	}
	return
}

// AnswerCallbackQuery send an answer to the given callback query
// https://core.telegram.org/bots/api#answercallbackquery
func (b *Bot) AnswerCallbackQuery(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("answerCallbackQuery", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsBool()
	}
	return
}

// SetMyCommands AnswerCallbackQuery send an answer to the given callback query
// https://core.telegram.org/bots/api#setmycommands
func (b *Bot) SetMyCommands(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("setMyCommands", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsBool()
	}
	return
}

// GetMyCommands send an answer to the given callback query
// https://core.telegram.org/bots/api#getmycommands
func (b *Bot) GetMyCommands() (result axon.A, err error) {
	var response interface{}
	if response, err = b.doGet("getMyCommands"); err == nil {
		result = response.([]interface{})
	}
	return
}

// EditMessageText edits a message text without resending
// https://core.telegram.org/bots/api#editemessagetext
func (b *Bot) EditMessageText(request axon.O) (result interface{}, err error) {
	var response interface{}
	if response, err = b.doPost("editMessageText", request); err == nil {
		v := axon.V{Value: response}
		if result, err = v.AsObject(); err != nil {
			result, err = v.AsBool()
		}
	}
	return
}

// EditMessageCaption edits a message caption without resending
// https://core.telegram.org/bots/api#editmessagecaption
func (b *Bot) EditMessageCaption(request axon.O) (result interface{}, err error) {
	var response interface{}
	if response, err = b.doPost("editMessageCaption", request); err == nil {
		v := axon.V{Value: response}
		if result, err = v.AsObject(); err != nil {
			result, err = v.AsBool()
		}
	}
	return
}

// EditMessageMedia edits a message media without resending
// https://core.telegram.org/bots/api#editmessagemedia
func (b *Bot) EditMessageMedia(request axon.O) (result interface{}, err error) {
	var response interface{}
	if response, err = b.doPost("editMessageMedia", request); err == nil {
		v := axon.V{Value: response}
		if result, err = v.AsObject(); err != nil {
			result, err = v.AsBool()
		}
	}
	return
}

// EditMessageReplyMarkup edits a message reply markup without resending
// https://core.telegram.org/bots/api#editmessagemedia
func (b *Bot) EditMessageReplyMarkup(request axon.O) (result interface{}, err error) {
	var response interface{}
	if response, err = b.doPost("editMessageReplyMarkup", request); err == nil {
		v := axon.V{Value: response}
		if result, err = v.AsObject(); err != nil {
			result, err = v.AsBool()
		}
	}
	return
}

// StopPoll stops a poll which was sent by the bot
// see https://core.telegram.org/bots/api#stoppoll
func (b *Bot) StopPoll(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("stopPoll", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// DeleteMessage deletes a message
// https://core.telegram.org/bots/api#deletemessage
func (b *Bot) DeleteMessage(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("deleteMessage", request); err == nil {
		v := axon.V{Value: response}
		result, err = v.AsBool()
	}
	return
}
