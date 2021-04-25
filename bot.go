// Package ubot provides types and
package ubot

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/sdurz/axon"
)

// matcherHandler encapsulates an UMatcher and the corresponding UHandler
type matcherHandler struct {
	matcher UMatcher
	handler UHandler
}

// evaluate execute the handler func if the matcher returns true
func (m *matcherHandler) evaluate(ctx context.Context, bot *Bot, message axon.O) (result bool, err error) {
	if m.matcher(bot, message) {
		result, err = m.handler(ctx, bot, message)
	}
	return
}

// Bot is the main type of ubot.
// It implements a bot API frontend.
type Bot struct {
	Configuration         Configuration
	apiClient             APIClient
	BotUser               UUser
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

// NewBot creates a new Bot for the given configuration
func NewBot(configuration *Configuration) (result *Bot) {
	if configuration.APIToken == "" {
		log.Fatalf("configuration has no APIToken")
	}
	if configuration.WorkerNo == 0 {
		configuration.WorkerNo = 5
	}
	result = &Bot{
		Configuration: *configuration,
		apiClient:     &httpApiClient{},
	}
	return
}

// methodURL transforms a method name in the corresponding API url
func (b *Bot) methodURL(method string) (result string) {
	if method == "" {
		panic("Emtpy method")
	}
	result = "https://api.telegram.org/bot" + b.Configuration.APIToken + "/" + method
	return
}

// getUpdatesSource
func (b *Bot) getUpdatesSource(ctx context.Context, updatesChan chan axon.O) {
	var nextUpdate int64 = 0
	var ok bool
	for {
		select {
		case <-ctx.Done():
			log.Println("done with getUpdatesSource")
			return
		default:
			getURL := b.methodURL("getUpdates") + "?offset=" + strconv.FormatInt(nextUpdate, 10)
			var responseUpdates interface{}
			responseUpdates, err := b.apiClient.GetJson(getURL)
			if err != nil {
				log.Println("Error while retrieving updates", err)
				continue
			}

			var updates axon.A
			if updates, ok = responseUpdates.([]interface{}); !ok {
				log.Fatalln("updates result not a JSON array")
			}

			if len(updates) > 0 {
				for _, update := range updates {
					var (
						updateID int64
						oUpdate  axon.O
					)
					if oUpdate, ok = update.(map[string]interface{}); !ok {
						log.Println("update not an axon.O")
						continue
					}
					if updateID, err = oUpdate.GetInteger("update_id"); err != nil {
						log.Println("update does not have an integer id")
						continue
					}
					if updateID > nextUpdate {
						nextUpdate = updateID
					}
					updatesChan <- oUpdate
				}
				nextUpdate++
			}
		}
	}
}

func (b *Bot) serverSource(ctx context.Context, updatesChan chan axon.O) {
	serverHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			body      []byte
			rawUpdate interface{}
			update    axon.O
			err       error
			ok        bool
		)

		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		log.Println(string(body))

		if err = json.Unmarshal(body, rawUpdate); err != nil {
			log.Printf("Error decoding body: %v", err)
			http.Error(w, "can't decode body", http.StatusBadRequest)
			return
		}

		if update, ok = rawUpdate.(map[string]interface{}); !ok {
			log.Printf("Error decoding body: %v", err)
			return
		}

		if err = b.process(ctx, update); err != nil {
			log.Println("Update processing error: ", err)
		}
	})

	if b.Configuration.WebhookUrl == "" {
		log.Fatal("empty webhook url")
	}

	mux := http.NewServeMux()
	mux.Handle("/bot/"+b.Configuration.APIToken, serverHandler)
	srv := &http.Server{
		Addr:    b.Configuration.ServerPort,
		Handler: mux,
	}

	if ok, err := b.SetWebhook(axon.O{"url": b.Configuration.WebhookUrl}); !ok || err != nil {
		log.Fatal("Can't set webhook")
		return
	}
	go http.ListenAndServe(b.Configuration.ServerPort, mux)
	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	log.Println("Server stopped")
}

func (b *Bot) AddChatMemberHandler(matcher UMatcher, handler UHandler) {
	b.chatMemberMHs = append(b.chatMemberMHs, matcherHandler{matcher: matcher, handler: handler})
}

func (b *Bot) Forever(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	source := b.serverSource
	if b.Configuration.LongPoll {
		source = b.getUpdatesSource
	}

	if user, err := b.GetMe(); err == nil {
		b.BotUser = *user
	} else {
		log.Fatal(err)
	}

	updates := make(chan axon.O)
	go source(ctx, updates)

	semaphore := make(chan int, b.Configuration.WorkerNo)
	for {
		select {
		case <-ctx.Done():
			log.Println("forever is over")
			return nil
		case update := <-updates:
			semaphore <- 1
			go func() {
				b.process(ctx, update)
				<-semaphore
			}()
		}
	}
}

func (b *Bot) process(ctx context.Context, update axon.O) (err error) {
	var (
		ok              bool
		stop            bool
		matcherHandlers []matcherHandler
		rawPayload      interface{}
		payload         axon.O
	)
	if rawPayload, ok = update["message"]; ok {
		matcherHandlers = b.messageMHs
	} else if rawPayload, ok = update["edited_message"]; ok {
		matcherHandlers = b.editedMessageMHs
	} else if rawPayload, ok = update["channel_post"]; ok {
		matcherHandlers = b.channelPostMHs
	} else if rawPayload, ok = update["edited_channel_post"]; ok {
		matcherHandlers = b.editedChannelPostMHs
	} else if rawPayload, ok = update["callback_query"]; ok {
		matcherHandlers = b.callbackQueryMHs
	} else if rawPayload, ok = update["inline_query"]; ok {
		matcherHandlers = b.inlineQueryMHs
	} else if rawPayload, ok = update["chosen_inline_result"]; ok {
		matcherHandlers = b.chosenInlineResultMHs
	} else {
		err = errors.New("update without data")
	}
	if payload, ok = rawPayload.(map[string]interface{}); !ok {
		err = errors.New("update payload not an axon.O")
		return
	}
	if err == nil {
		for _, ha := range matcherHandlers {
			if stop, err = ha.evaluate(ctx, b, payload); err != nil {
				log.Println(err)
				break
			}
			if stop {
				break
			}
		}
	}
	return
}

func (b *Bot) doGet(method string, data axon.O) (interface{}, error) {
	return b.apiClient.GetJson(b.methodURL(method))
}

func (b *Bot) doPost(method string, request axon.O) (interface{}, error) {
	return b.apiClient.PostJson(b.methodURL(method), request)
}

func (b *Bot) doPostMultipart(method string, request axon.O) (interface{}, error) {
	return b.apiClient.PostMultipart(b.methodURL(method), request)
}

func (b *Bot) AddMessageHandler(matcher UMatcher, handler UHandler) {
	b.messageMHs = append(b.messageMHs, matcherHandler{matcher: matcher, handler: handler})
}

func (b *Bot) AddChannelPostHandler(matcher UMatcher, handler UHandler) {
	b.channelPostMHs = append(b.channelPostMHs, matcherHandler{matcher: matcher, handler: handler})
}

func (b *Bot) AddEditedMessageHandler(matcher UMatcher, handler UHandler) {
	b.editedMessageMHs = append(b.editedMessageMHs, matcherHandler{matcher: matcher, handler: handler})
}

func (b *Bot) AddEditedChannelPostHandler(matcher UMatcher, handler UHandler) {
	b.editedChannelPostMHs = append(b.editedChannelPostMHs, matcherHandler{matcher: matcher, handler: handler})
}

func (b *Bot) AddCallbackQueryHandler(matcher UMatcher, handler UHandler) {
	b.callbackQueryMHs = append(b.callbackQueryMHs, matcherHandler{matcher: matcher, handler: handler})
}

func (b *Bot) AddMyChatMemberHandler(matcher UMatcher, handler UHandler) {
	b.myChatMemberMHs = append(b.myChatMemberMHs, matcherHandler{matcher: matcher, handler: handler})
}

// API Mappings

// GetMe returns basic information about the bot in form of a User object.
// see https://core.telegram.org/bots/api#getme
func (b *Bot) GetMe() (result *UUser, err error) {
	var (
		uResult UUser
		iResult interface{}
		oResult axon.O
		ok      bool
	)
	if iResult, err = b.doGet("getMe", nil); err != nil {
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
	_, err = b.doGet("logOut", nil)
	return
}

// Close closea the bot instance
// see https://core.telegram.org/bots/api#close
func (b *Bot) Close() (err error) {
	_, err = b.doGet("close", nil)
	return
}

// SendMessage sends a text message
// see https://core.telegram.org/bots/api#sendmessage
func (b *Bot) SendMessage(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("sendMessage", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// ForwardMessage forwards messages of any kind
// see https://core.telegram.org/bots/api#forwardmessage
func (b *Bot) ForwardMessage(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("forwardMessage", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// CopyMessage copy messages of any kind. The method is analogous to the method forwardMessage, but the copied message doesn't have a link to the original message.
// see https://core.telegram.org/bots/api#copymessage
func (b *Bot) CopyMessage(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("copyMessage", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendPhoto sends a photo
// see https://core.telegram.org/bots/api#sendphoto
func (b *Bot) SendPhoto(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendPhoto", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendAudio sends an audio
// see https://core.telegram.org/bots/api#sendaudio
func (b *Bot) SendAudio(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendAudio", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendVideo sends a video
// see https://core.telegram.org/bots/api#senddocument
func (b *Bot) SendDocument(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendDocument", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendVideo sends a video
// see https://core.telegram.org/bots/api#sendvideo
func (b *Bot) SendVideo(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendVideo", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendAnimation sends an animation
// see https://core.telegram.org/bots/api#sendanimation
func (b *Bot) SendAnimation(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendAnimation", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendVoice sends a voice
// see https://core.telegram.org/bots/api#sendvoice
func (b *Bot) SendVoice(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendVoice", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendVoice sends a video note
// see https://core.telegram.org/bots/api#sendvideonote
func (b *Bot) SendVideoNote(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendVideoNote", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendVoice sends a media group
// see https://core.telegram.org/bots/api#sendmediagroup
func (b *Bot) SendMediaGroup(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendMediaGroup", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendLocation sends a location
// see https://core.telegram.org/bots/api#sendlocation
func (b *Bot) SendLocation(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("sendLocation", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// EditMessageLiveLocation sends a location
// see https://core.telegram.org/bots/api#editmessagelivelocation
func (b *Bot) EditMessageLiveLocation(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("editMessageLiveLocation", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// StopMessageLiveLocation sends a location
// see https://core.telegram.org/bots/api#stopmessagelivelocation
func (b *Bot) StopMessageLiveLocation(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("stopMessageLiveLocation", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendVenue sends a venue
// see https://core.telegram.org/bots/api#sendvenue
func (b *Bot) SendVenue(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("sendVenue", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendContact sends a venue
// see https://core.telegram.org/bots/api#sendcontact
func (b *Bot) SendContact(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("sendContact", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendPoll sends a poll
// see https://core.telegram.org/bots/api#sendpoll
func (b *Bot) SendPoll(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("sendPoll", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendDice sends a dice
// see https://core.telegram.org/bots/api#senddice
func (b *Bot) SendDice(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("sendDice", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendChatAction sends a chat action
// see https://core.telegram.org/bots/api#sendchataction
func (b *Bot) SendChatAction(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("sendChatAction", request); err == nil {
		result = response.(bool)
	}
	return
}

// GetUserProfilePhotos gets user profile photos.
// see https://core.telegram.org/bots/api#getuserprofilephotos
func (b *Bot) GetUserProfilePhotos(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("getUserProfilesPhotos", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// GetFile gets basic info about a file and prepare it for downloading.
// see https://core.telegram.org/bots/api#getfile
func (b *Bot) GetFile(fileId string) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doGet("getFile?file_id="+fileId, nil); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// KickChatMember kicks a user from a group, a supergroup or a channel.
// see https://core.telegram.org/bots/api#kickchatmember
func (b *Bot) KickChatMember(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("kickChatMember", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// UnbanChatMember unban a previously kicked user in a supergroup or channel.
// see https://core.telegram.org/bots/api#unbanchatmember
func (b *Bot) UnbanChatMember(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("unbanChatMember", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// RestrictChatMember unban a previously kicked user in a supergroup or channel.
// see https://core.telegram.org/bots/api#restrictchatmember
func (b *Bot) RestrictChatMember(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("restrictChatMember", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// PromoteChatMember unban a previously kicked user in a supergroup or channel.
// see https://core.telegram.org/bots/api#promotechatmember
func (b *Bot) PromoteChatMember(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("promoteChatMember", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SetWebhook implements setWebhook from Telegram Bot API.
// see https://core.telegram.org/bots/api#setwebhook
func (b *Bot) SetWebhook(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("setWebhook", request); err == nil {
		result = response.(bool)
	}
	return
}

// DeleteWebhook implements deleteWebhook from Telegram Bot API
// see https://core.telegram.org/bots/api#deletewebhook
func (b *Bot) DeleteWebhook(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("deleteWebhook", request); err == nil {
		result = response.(bool)
	}
	return
}

// GetWebhookInfo get current webhook status
// https://core.telegram.org/bots/api#getwebhookinfo
func (b *Bot) GetWebhookInfo() (result axon.O, err error) {
	var response interface{}
	if response, err = b.doGet("getWebhookInfo", result); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

func (b *Bot) AnswerCallbackQuery(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("answerCallbackQuery", request); err == nil {
		result = response.(bool)
	}
	return
}

func (b *Bot) PinChatMessage(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("answerCallbackQuery", request); err == nil {
		result = response.(bool)
	}
	return
}

func (b *Bot) UnpinChatMessage(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("answerCallbackQuery", request); err == nil {
		result = response.(bool)
	}
	return
}

func (b *Bot) UnpinAllChatMessages(request axon.O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("answerCallbackQuery", request); err == nil {
		result = response.(bool)
	}
	return
}

func (b *Bot) GetChat(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("getChat", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

func (b *Bot) LeaveChat(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("leaveChat", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

func (b *Bot) GetChatMembersCount(request axon.O) (result int64, err error) {
	var response interface{}
	if response, err = b.doPost("getChat", request); err == nil {
		result = response.(int64)
	}
	return
}

func (b *Bot) SetMyCommands(request axon.O) (result axon.O, err error) {
	var response interface{}
	if response, err = b.doPost("setMyCommands", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

func (b *Bot) GetMyCommands() (result axon.O, err error) {
	var response interface{}
	if response, err = b.doGet("getMyCommands", result); err == nil {
		result = response.(map[string]interface{})
	}
	return
}
