# Elide
**Elide** is a multi-protocol extensional API to provide support for some missing methods from the official Telegram Bot API like resolveUsername, getMessages etc.

**Note**: `It's work-in-progess at the moment.`

## Supported Methods

### endpoint`/resolveUsername`
This method can be used to resolve a telegram username.

Parameters:
- `username` (type string): Username of the telegram channel/chat. *

### endpoint`/getMessages`
This method can be used to get list of telegram messages in a chat.

Parameters:
- `chat_id` (type int64): Unique identifier of the telegram channel/chat. *
- `message_ids` (type []int64): List of ids of messages to be fetched. *

### endpoint`/deleteMessages`
This method can be used to delete list of telegram messages in a chat.

Parameters:
- `chat_id` (type int64): Unique identifier of the telegram channel/chat. (optional)
- `message_ids` (type []int64): List of ids of messages to be fetched. *
- `revoke` (type bool): Pass `true` to revoke messages from a private chat. (optional)


## License
Licensed under **GNU AFFERO GENERAL PUBLIC LICENSE V3**
