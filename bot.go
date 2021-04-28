// Package ubot provides types and
package ubot

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/sdurz/axon"
)

// Configuration struct holds configuration data for the bot
type Configuration struct {
	APIToken   string `json:"api_token"`
	ServerPort string `json:"server_port"`
	WebhookUrl string `json:"webhook_url"`
	WorkerNo   int    `json:"worker_no"`
}

// Bot is the main type of ubot.
// It implements a bot API frontend.
type Bot struct {
	Configuration         Configuration
	apiClient             apiClient
	BotUser               User
	messageMHs            []matcherHandler
	editedMessageMHs      []matcherHandler
	channelPostMHs        []matcherHandler
	editedChannelPostMHs  []matcherHandler
	inlineQueryMHs        []matcherHandler
	chosenInlineResultMHs []matcherHandler
	callbackQueryMHs      []matcherHandler
	shippingQueryMHs      []matcherHandler
	preCheckoutQueryMHs   []matcherHandler
	pollMHs               []matcherHandler
	pollAnswerMHs         []matcherHandler
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

// matcherHandler encapsulates an Matcher and the corresponding Handler
type matcherHandler struct {
	matcher Matcher
	handler Handler
}

// evaluate execute the handler func if the matcher returns true
func (m *matcherHandler) evaluate(ctx context.Context, bot *Bot, message axon.O) (result bool, err error) {
	if m.matcher(bot, message) {
		result, err = m.handler(ctx, bot, message)
	}
	return
}

// AddMessageHandler adds an handler for message updates.
func (b *Bot) AddMessageHandler(matcher Matcher, handler Handler) {
	b.messageMHs = append(b.messageMHs, matcherHandler{matcher: matcher, handler: handler})
}

// AddEditedMessageHandler adds an handler for edited_message updates.
func (b *Bot) AddEditedMessageHandler(matcher Matcher, handler Handler) {
	b.editedMessageMHs = append(b.editedMessageMHs, matcherHandler{matcher: matcher, handler: handler})
}

// AddChannelPostHandler adds an handler for channel_post updates.
func (b *Bot) AddChannelPostHandler(matcher Matcher, handler Handler) {
	b.channelPostMHs = append(b.channelPostMHs, matcherHandler{matcher: matcher, handler: handler})
}

// AddEditedChannelPostHandler adds an handler for edited_channel_post updates.
func (b *Bot) AddEditedChannelPostHandler(matcher Matcher, handler Handler) {
	b.editedChannelPostMHs = append(b.editedChannelPostMHs, matcherHandler{matcher: matcher, handler: handler})
}

// AddInlineQueryHandler adds an handler for inline_query updates.
func (b *Bot) AddInlineQueryHandler(matcher Matcher, handler Handler) {
	b.inlineQueryMHs = append(b.inlineQueryMHs, matcherHandler{matcher: matcher, handler: handler})
}

// AddChosenInlineResultHandler adds an handler for inline_query updates.
func (b *Bot) AddChosenInlineResultHandler(matcher Matcher, handler Handler) {
	b.chosenInlineResultMHs = append(b.chosenInlineResultMHs, matcherHandler{matcher: matcher, handler: handler})
}

// AddCallbackQueryHandler adds an handler for callback_query updates.
func (b *Bot) AddCallbackQueryHandler(matcher Matcher, handler Handler) {
	b.callbackQueryMHs = append(b.callbackQueryMHs, matcherHandler{matcher: matcher, handler: handler})
}

// AddShippingQueryHandler adds an handler for callback_query updates.
func (b *Bot) AddShippingQueryHandler(matcher Matcher, handler Handler) {
	b.shippingQueryMHs = append(b.shippingQueryMHs, matcherHandler{matcher: matcher, handler: handler})
}

// AddPreCheckoutQueryHandler adds an handler for callback_query updates.
func (b *Bot) AddPreCheckoutQueryHandler(matcher Matcher, handler Handler) {
	b.preCheckoutQueryMHs = append(b.preCheckoutQueryMHs, matcherHandler{matcher: matcher, handler: handler})
}

// AddPollHandler adds an handler for callback_query updates.
func (b *Bot) AddPollHandler(matcher Matcher, handler Handler) {
	b.pollMHs = append(b.pollMHs, matcherHandler{matcher: matcher, handler: handler})
}

// AddPollAnswerHandler adds an handler for callback_query updates.
func (b *Bot) AddPollAnswerHandler(matcher Matcher, handler Handler) {
	b.pollAnswerMHs = append(b.pollAnswerMHs, matcherHandler{matcher: matcher, handler: handler})
}

// AddMyChatMemberHandler adds an handler for my_chat_member updates.
func (b *Bot) AddMyChatMemberHandler(matcher Matcher, handler Handler) {
	b.myChatMemberMHs = append(b.myChatMemberMHs, matcherHandler{matcher: matcher, handler: handler})
}

// AddChatMemberHandler adds an handler for chat_member updates.
func (b *Bot) AddChatMemberHandler(matcher Matcher, handler Handler) {
	b.chatMemberMHs = append(b.chatMemberMHs, matcherHandler{matcher: matcher, handler: handler})
}

// Forever starts the bot and processes updates until context is done.
func (b *Bot) Forever(ctx context.Context, wg *sync.WaitGroup, source UpdatesSource) error {
	defer wg.Done()

	if user, err := b.GetMe(); err == nil {
		b.BotUser = *user
	} else {
		log.Fatal(err)
	}

	updates := make(chan axon.O)
	go source(b, ctx, updates)

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

// methodURL transforms a method name in the corresponding API url
func (b *Bot) methodURL(method string) (result string) {
	if method == "" {
		panic("Emtpy method")
	}
	result = "https://api.telegram.org/bot" + b.Configuration.APIToken + "/" + method
	return
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
	} else if rawPayload, ok = update["inline_query"]; ok {
		matcherHandlers = b.inlineQueryMHs
	} else if rawPayload, ok = update["chosen_inline_result"]; ok {
		matcherHandlers = b.chosenInlineResultMHs
	} else if rawPayload, ok = update["callback_query"]; ok {
		matcherHandlers = b.callbackQueryMHs
	} else if rawPayload, ok = update["shipping_query"]; ok {
		matcherHandlers = b.shippingQueryMHs
	} else if rawPayload, ok = update["pre_checkout_query"]; ok {
		matcherHandlers = b.preCheckoutQueryMHs
	} else if rawPayload, ok = update["poll"]; ok {
		matcherHandlers = b.pollMHs
	} else if rawPayload, ok = update["poll_answer"]; ok {
		matcherHandlers = b.pollAnswerMHs
	} else if rawPayload, ok = update["my_chat_member"]; ok {
		matcherHandlers = b.myChatMemberMHs
	} else if rawPayload, ok = update["chat_member"]; ok {
		matcherHandlers = b.chatMemberMHs
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

func (b *Bot) doGet(method string) (interface{}, error) {
	return b.apiClient.GetJson(b.methodURL(method))
}

func (b *Bot) doPost(method string, request axon.O) (interface{}, error) {
	return b.apiClient.PostJson(b.methodURL(method), request)
}

func (b *Bot) doPostMultipart(method string, request axon.O) (interface{}, error) {
	return b.apiClient.PostMultipart(b.methodURL(method), request)
}
