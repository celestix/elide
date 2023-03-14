package elide

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/gotd/td/tg"
)

type Client struct {
	Port, Protocol int
	conn           net.Conn
	addr           string
}

type Response[T comparable] struct {
	Ok     bool   `json:"ok"`
	Error  string `json:"error,omitempty"`
	Result T      `json:"result,omitempty"`
}

func NewClient(protocol, port int) (*Client, error) {
	c := &Client{
		Protocol: protocol,
		Port:     port,
	}
	switch protocol {
	case 0:
		c.addr = fmt.Sprintf("http://0.0.0.0:%d", port)
	case 1:
		c.addr = fmt.Sprintf("0.0.0.0:%d", port)
		conn, err := net.Dial("tcp", c.addr)
		if err != nil {
			return nil, err
		}
		c.conn = conn
	}
	return c, nil
}

type Body struct {
	Method string `json:"method"`
	Data   any    `json:"data"`
}

type DeleteMessagesBody struct {
	ChatId     int64 `json:"chat_id,omitempty"`
	MessageIds []int `json:"message_ids"`
	Revoke     bool  `json:"revoke,omitempty"`
}

type GetMessagesBody struct {
	ChatId     int64 `json:"chat_id"`
	MessageIds []int `json:"message_ids"`
}

type ResolveUsernameBody struct {
	Username string `json:"username"`
}

type GetChatPhotoBody struct {
	// User ID
	ChatId any `json:"chat_id"`
}

type GetProfilePhotosBody struct {
	// User ID
	UserID any `json:"user_id"`
	// Number of list elements to be skipped
	Offset int `json:"offset,omitempty"`
	// If a positive value was transferred, the method will return only photos with IDs less
	// than the set one
	MaxID int64 `json:"max_id,omitempty"`
	// Number of list elements to be returned
	Limit int `json:"limit,omitempty"`
}

// Message represents TL type `message#38116ee0`.
// A message
//
// See https://core.telegram.org/constructor/message for reference.
type Message struct {
	// Is this an outgoing message
	Out bool
	// Whether we were mentioned¹ in this message
	//
	// Links:
	//  1) https://core.telegram.org/api/mentions
	Mentioned bool
	// Whether there are unread media attachments in this message
	MediaUnread bool
	// Whether this is a silent message (no notification triggered)
	Silent bool
	// Whether this is a channel post
	Post bool
	// Whether this is a scheduled message¹
	//
	// Links:
	//  1) https://core.telegram.org/api/scheduled-messages
	FromScheduled bool
	// This is a legacy message: it has to be refetched with the new layer
	Legacy bool
	// Whether the message should be shown as not modified to the user, even if an edit date
	// is present
	EditHide bool
	// Whether this message is pinned¹
	//
	// Links:
	//  1) https://core.telegram.org/api/pin
	Pinned bool
	// Whether this message is protected¹ and thus cannot be forwarded
	//
	// Links:
	//  1) https://telegram.org/blog/protected-content-delete-by-date-and-more
	Noforwards bool
	// ID of the message
	ID int
	// ID of the sender of the message
	//
	// Use SetFromID and GetFromID helpers.
	FromID Peer
	// Peer ID, the chat where this message was sent
	PeerID Peer
	// Info about forwarded messages
	FwdFrom MessageFwdHeader
	// ID of the inline bot that generated the message
	ViaBotID int64
	// Reply information
	ReplyTo MessageReplyHeader
	// Date of the message
	Date int
	// The message
	Message string
	// Media attachment
	// TODO: Create a proper type for it
	Media json.RawMessage
	// Reply markup (bot/inline keyboards)
	// TODO: Create a proper type for it
	ReplyMarkup json.RawMessage
	// Message entities¹ for styled text
	//
	// Links:
	//  1) https://core.telegram.org/api/entities
	// TODO: Create a proper type for it
	Entities []json.RawMessage
	// View count for channel posts
	Views int
	// Forward counter
	Forwards int
	// Info about post comments (for channels) or message replies (for groups)¹
	//
	// Links:
	//  1) https://core.telegram.org/api/threads
	Replies tg.MessageReplies
	// Last edit date of this message
	EditDate int
	// Name of the author of this message for channel posts (with signatures enabled)
	PostAuthor string
	// Multiple media messages sent using messages.sendMultiMedia¹ with the same grouped ID
	// indicate an album or media group²
	//
	// Links:
	//  1) https://core.telegram.org/method/messages.sendMultiMedia
	//  2) https://core.telegram.org/api/files#albums-grouped-media
	GroupedID int64
	// Contains the reason why access to this message must be restricted.
	RestrictionReason []RestrictionReason
	// Time To Live of the message, once message.date+message.ttl_period === time(), the
	// message will be deleted on the server, and must be deleted locally as well.
	TTLPeriod int
}

// MessageReplyHeader represents TL type `messageReplyHeader#a6d57763`.
// Message replies and thread¹ information
//
// Links:
//  1. https://core.telegram.org/api/threads
//
// See https://core.telegram.org/constructor/messageReplyHeader for reference.
type MessageReplyHeader struct {
	// Whether this message replies to a scheduled message
	ReplyToScheduled bool
	// ID of message to which this message is replying
	ReplyToMsgID int
	// For replies sent in channel discussion threads¹ of which the current user is not a
	// member, the discussion group ID
	//
	// Links:
	//  1) https://core.telegram.org/api/threads
	ReplyToPeerID Peer
	// ID of the message that started this message thread¹
	//
	// Links:
	//  1) https://core.telegram.org/api/threads
	ReplyToTopID int
}

// MessageFwdHeader represents TL type `messageFwdHeader#5f777dce`.
// Info about a forwarded message
//
// See https://core.telegram.org/constructor/messageFwdHeader for reference.
type MessageFwdHeader struct {
	// Whether this message was imported from a foreign chat service, click here for more
	// info »¹
	//
	// Links:
	//  1) https://core.telegram.org/api/import
	Imported bool
	// The ID of the user that originally sent the message
	FromID Peer
	// The name of the user that originally sent the message
	FromName string
	// When was the message originally sent
	Date int
	// ID of the channel message that was forwarded
	ChannelPost int
	// For channels and if signatures are enabled, author of the channel message
	PostAuthor string
	// Only for messages forwarded to the current user (inputPeerSelf), full info about the
	// user/channel that originally sent the message
	SavedFromPeer Peer
	// Only for messages forwarded to the current user (inputPeerSelf), ID of the message
	// that was forwarded from the original user/channel
	SavedFromMsgID int
	// PSA type
	PsaType string
}

// MessageReplies represents TL type `messageReplies#83d60fc2`.
// Info about the comment section of a channel post, or a simple message thread¹
//
// Links:
//  1. https://core.telegram.org/api/threads
//
// See https://core.telegram.org/constructor/messageReplies for reference.
type MessageReplies struct {
	// Whether this constructor contains information about the comment section of a channel
	// post, or a simple message thread¹
	//
	// Links:
	//  1) https://core.telegram.org/api/threads
	Comments bool
	// Contains the total number of replies in this thread or comment section.
	Replies int
	// For channel post comments, contains information about the last few comment posters for
	// a specific thread, to show a small list of commenter profile pictures in client
	// previews.
	//
	// Use SetRecentRepliers and GetRecentRepliers helpers.
	RecentRepliers []Peer
	// For channel post comments, contains the ID of the associated discussion supergroup¹
	//
	// Links:
	//  1) https://core.telegram.org/api/discussion
	ChannelID int64
	// ID of the latest message in this thread or comment section.
	MaxID int
	// Contains the ID of the latest read message in this thread or comment section.
	ReadMaxID int
}

type Peer struct {
	UserId, ChatId, ChannelId int64
}

// RestrictionReason represents TL type `restrictionReason#d072acb4`.
// Restriction reason.
// Contains the reason why access to a certain object must be restricted. Clients are
// supposed to deny access to the channel if the platform field is equal to all or to the
// current platform (ios, android, wp, etc.). Platforms can be concatenated (ios-android,
// ios-wp), unknown platforms are to be ignored. The text is the error message that
// should be shown to the user.
//
// See https://core.telegram.org/constructor/restrictionReason for reference.
type RestrictionReason struct {
	// Platform identifier (ios, android, wp, all, etc.), can be concatenated with a dash as
	// separator (android-ios, ios-wp, etc)
	Platform string
	// Restriction reason (porno, terms, etc.)
	Reason string
	// Error message to be shown to the user
	Text string
}

type Chat struct {
	// ID of the chat
	ID int64
	// Title
	Title string
	// Whether this chat indicates the currently logged in bot
	Self bool
	// Whether the account of this user was deleted
	Deleted bool
	// Is this user a bot?
	Bot bool
	// Can the bot see all messages in groups?
	BotChatHistory bool
	// Can the bot be added to groups?
	BotNochats bool
	// Whether this user is verified
	Verified bool
	// Access to this user must be restricted for the reason specified in restriction_reason
	Restricted bool
	// Whether the bot can request our geolocation in inline mode
	BotInlineGeo bool
	// Whether this is an official support user
	Support bool
	// This may be a scam user
	Scam bool
	// If set, the profile picture for this user should be refetched
	ApplyMinPhoto bool
	// If set, this user was reported by many users as a fake or scam user: be careful when
	// interacting with them.
	Fake bool
	//
	BotAttachMenu bool
	// Whether this user is a Telegram Premium user
	Premium bool
	//
	AttachMenuEnabled bool
	// First name
	FirstName string
	// Last name
	LastName string
	// Phone number
	Phone string
	// Version of the bot_info field in userFull¹, incremented every time it changes
	//
	// Links:
	//  1) https://core.telegram.org/constructor/userFull
	BotInfoVersion int
	// Inline placeholder for this inline bot
	BotInlinePlaceholder string
	// Language code of the user
	LangCode string
	// Is this a channel?
	Broadcast bool
	// Is this a supergroup?
	Megagroup bool
	// Whether signatures are enabled (channels)
	Signatures bool
	// Whether this channel has a private join link
	HasLink bool
	// Whether this chanel has a geoposition
	HasGeo bool
	// Whether slow mode is enabled for groups to prevent flood in chat
	SlowmodeEnabled bool
	// Whether this supergroup¹ is a gigagroup
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	Gigagroup bool
	// Whether a user needs to join the supergroup before they can send messages: can be
	// false only for discussion groups »¹, toggle using channels.toggleJoinToSend²
	//
	// Links:
	//  1) https://core.telegram.org/api/discussion
	//  2) https://core.telegram.org/method/channels.toggleJoinToSend
	JoinToSend bool
	// Whether a user's join request will have to be approved by administrators¹, toggle
	// using channels.toggleJoinToSend²
	//
	// Links:
	//  1) https://core.telegram.org/api/invites#join-requests
	//  2) https://core.telegram.org/method/channels.toggleJoinRequest
	JoinRequest bool
	// Username
	Username string
	// Whether the current user is the creator of the group
	Creator bool
	// Whether the current user has left the group
	Left bool
	// Whether the group was migrated¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	Deactivated bool
	// Whether a group call is currently active
	CallActive bool
	// Whether there's anyone in the group call
	CallNotEmpty bool
	// Whether this group is protected¹, thus does not allow forwarding messages from it
	//
	// Links:
	//  1) https://telegram.org/blog/protected-content-delete-by-date-and-more
	Noforwards bool
	// Chat photo
	Photo tg.ChatPhotoEmpty
	// Participant count
	ParticipantsCount int
	// Date of creation of the group
	Date int
	// Used in basic groups to reorder updates and make sure that all of them were received.
	Version int
	// Means this chat was upgraded¹ to a supergroup
	MigratedTo struct {
		// Channel ID
		ChannelID int64
		// Access hash taken from the channel¹ constructor
		//
		// Links:
		//  1) https://core.telegram.org/constructor/channel
		AccessHash int64
	}
	// Admin rights¹ of the user in the group
	//
	// Links:
	//  1) https://core.telegram.org/api/rights
	AdminRights ChatAdminRights
	// Default banned rights¹ of all users in the group
	//
	// Links:
	//  1) https://core.telegram.org/api/rights
	DefaultBannedRights ChatBannedRights
}

// ChatAdminRights represents TL type `chatAdminRights#5fb224d5`.
// Represents the rights of an admin in a channel/supergroup¹.
//
// Links:
//  1. https://core.telegram.org/api/channel
//
// See https://core.telegram.org/constructor/chatAdminRights for reference.
type ChatAdminRights struct {
	// If set, allows the admin to modify the description of the channel/supergroup¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	ChangeInfo bool
	// If set, allows the admin to post messages in the channel¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	PostMessages bool
	// If set, allows the admin to also edit messages from other admins in the channel¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	EditMessages bool
	// If set, allows the admin to also delete messages from other admins in the channel¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	DeleteMessages bool
	// If set, allows the admin to ban users from the channel/supergroup¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	BanUsers bool
	// If set, allows the admin to invite users in the channel/supergroup¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	InviteUsers bool
	// If set, allows the admin to pin messages in the channel/supergroup¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	PinMessages bool
	// If set, allows the admin to add other admins with the same (or more limited)
	// permissions in the channel/supergroup¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	AddAdmins bool
	// Whether this admin is anonymous
	Anonymous bool
	// If set, allows the admin to change group call/livestream settings
	ManageCall bool
	// Set this flag if none of the other flags are set, but you still want the user to be an
	// admin.
	Other bool
}

// ChatBannedRights represents TL type `chatBannedRights#9f120418`.
// Represents the rights of a normal user in a supergroup/channel/chat¹. In this case,
// the flags are inverted: if set, a flag does not allow a user to do X.
//
// Links:
//  1. https://core.telegram.org/api/channel
//
// See https://core.telegram.org/constructor/chatBannedRights for reference.
type ChatBannedRights struct {
	// If set, does not allow a user to view messages in a supergroup/channel/chat¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	ViewMessages bool
	// If set, does not allow a user to send messages in a supergroup/chat¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	SendMessages bool
	// If set, does not allow a user to send any media in a supergroup/chat¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	SendMedia bool
	// If set, does not allow a user to send stickers in a supergroup/chat¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	SendStickers bool
	// If set, does not allow a user to send gifs in a supergroup/chat¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	SendGifs bool
	// If set, does not allow a user to send games in a supergroup/chat¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	SendGames bool
	// If set, does not allow a user to use inline bots in a supergroup/chat¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	SendInline bool
	// If set, does not allow a user to embed links in the messages of a supergroup/chat¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	EmbedLinks bool
	// If set, does not allow a user to send polls in a supergroup/chat¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	SendPolls bool
	// If set, does not allow any user to change the description of a supergroup/chat¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	ChangeInfo bool
	// If set, does not allow any user to invite users in a supergroup/chat¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	InviteUsers bool
	// If set, does not allow any user to pin messages in a supergroup/chat¹
	//
	// Links:
	//  1) https://core.telegram.org/api/channel
	PinMessages bool
	// Validity of said permissions (it is considered forever any value less then 30 seconds
	// or more then 366 days).
	UntilDate int
}

// ChatPhoto represents TL type `chatPhoto#1c6e1c11`.
// Group profile photo.
//
// See https://core.telegram.org/constructor/chatPhoto for reference.
type ChatPhoto struct {
	// Whether the user has an animated profile picture
	HasVideo bool `json:"HasVideo,omitempty"`
	// Photo ID
	PhotoID int64 `json:"PhotoID,omitempty"`
	// Stripped thumbnail¹
	//
	// Links:
	//  1) https://core.telegram.org/api/files#stripped-thumbnails
	StrippedThumb []byte `json:"StrippedThumb,omitempty"`
	// DC where this photo is stored
	DCID int `json:"DCID,omitempty"`
}
