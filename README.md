# Elide
**Elide** is a multi-protocol extensional API to provide support for some missing methods from the official Telegram Bot API like resolveUsername, getMessages etc.

**Note**: `It's work-in-progess at the moment.`

## Installation

- You can install it directly into your machine using go install command:
  ```go install github.com/anonyindian/elide```
- **Manually Compiling** - You can also compile it manually without any extra steps using go build command:
  ```go build .```
  
## Usage

Once you install elide, go through the following process to run it:

1. Run `elide init` command to generate the config file in the current directory.
2. Fill up the details in the generated config file correctly.
3. Run `elide run` in the current directory to start the server.

Once you follow the above steps, your elide instance shall start successfully.
You can reach me out on [telegram](https://t.me/CaptainPicard) in case you need any help regarding it, etc.


## Supported Methods

### > endpoint`/resolveUsername`
This method can be used to resolve a telegram username.

**Parameters**:
- `username` (type string) (required): Username of the telegram channel/chat.

### > endpoint`/getMessages`
This method can be used to get list of telegram messages in a chat.

**Parameters**:
- `chat_id` (type int64) (required): Unique identifier of the telegram channel/chat.
- `message_ids` (type []int64) (required): List of ids of messages to be fetched. 

### > endpoint`/deleteMessages`
This method can be used to delete list of telegram messages in a chat.

**Parameters**:
- `chat_id` (type int64) (optional): Unique identifier of the telegram channel/chat. 
- `message_ids` (type []int64) (required): List of ids of messages to be fetched.
- `revoke` (type bool) (optional): Pass `true` to revoke messages from a private chat.


## License
Licensed under **GNU AFFERO GENERAL PUBLIC LICENSE V3**
