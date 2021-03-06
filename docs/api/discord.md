# Discord
- Type = `discords`

> **Note**: You should create an env file in /.
```dotenv
DISCORD_BOT_SECRET=<DISCORD_SECRET>
```

### Connection
```json
{
    "type": "discord",
    "data": {
        "guildId": "844018666315711474"
    }
}
```

### SubTypes & Data
Adding discord node to the scenario.
- `sendMessage` (Send a message in a channel)
    - Data request to create:
    
        ```json
        {
            "data": {
              "channelId": "128373",
              "content": "this is the message!"
            }
        }
        ```
    
    - Data response
        ```json
        {
            "data": {
              "channelId": "128373",
              "content": "this is the message!"
            }
        }
        ```
      
- `pinMessage` (Pin a message in a channel)
    - Data request to create:
    
        ```json
        {
            "data": {
              "channelId": "128373",
              "messageId": "31231234"
            }
        }
        ```
    
    - Data response
        ```json
        {
            "data": {
              "channelId": "128373",
              "messageId": "31231234"
            }
        }
        ```

- `reactMessage` (react to a message in a channel)
    > **Note:** The emoji should be just like this. Cure about the format. 
    - Data request to create:
    
        ```json
        {
            "data": {
              "channelId": "128373",
              "messageId": "31231234",
              "react": "😁"
            }
        }
        ```
    
    - Data response
        ```json
        {
            "data": {
              "channelId": "128373",
              "messageId": "31231234"
            }
        }
        ```

### /api/services/discord/guilds/:guildId/channels/text `GET`
Returns the text channels of a guild.
- Response
    ```json
    {
        "channels": [
            {
                "id": "844018666315710479",
                "guild_id": "844018666315710476",
                "name": "general",
                "topic": "",
                "type": 0,
                "last_message_id": "847540869434703932",
                "nsfw": false,
                "icon": "",
                "position": 0,
                "bitrate": 0,
                "user_limit": 0,
                "parent_id": "844018666315710477",
                "rate_limit_per_user": 0,
                "owner_id": "",
                "application_id": ""
            },
            {
                "id": "845633435011514369",
                "guild_id": "844018666315710476",
                "name": "groovy",
                "topic": "",
                "type": 0,
                "last_message_id": "846502029579780117",
                "nsfw": false,
                "icon": "",
                "position": 2,
                "bitrate": 0,
                "user_limit": 0,
                "parent_id": "844018666315710477",
                "rate_limit_per_user": 0,
                "owner_id": "",
                "application_id": ""
            },
            {
                "id": "846051939149807666",
                "guild_id": "844018666315710476",
                "name": "temp",
                "topic": "",
                "type": 0,
                "last_message_id": "847540868630315049",
                "nsfw": false,
                "icon": "",
                "position": 3,
                "bitrate": 0,
                "user_limit": 0,
                "parent_id": "844018666315710477",
                "rate_limit_per_user": 0,
                "owner_id": "",
                "application_id": ""
            }
        ]
    }
    ```
  
 ### /api/services/discord/channels/:channelId/messages `GET`
 Returns the messages of a single text channels.
 - Response
     ```json
     {
         "messages": [
            {
                "id": "847626806651519048",
                "type": 6,
                "content": "that is it",
                "channel_id": "846160866000371753",
                "attachments": [],
                "embeds": [],
                "mentions": [],
                "mention_roles": [],
                "pinned": false,
                "mention_everyone": false,
                "tts": false,
                "timestamp": "2021-05-28T00:06:18.850000+00:00",
                "edited_timestamp": null,
                "flags": 0,
                "components": []
            },
            {
                "id": "847624201682812958",
                "type": 0,
                "content": "hello",
                "channel_id": "846160866000371753",
                "attachments": [],
                "embeds": [],
                "mentions": [],
                "mention_roles": [],
                "pinned": true,
                "mention_everyone": false,
                "tts": false,
                "timestamp": "2021-05-27T23:55:57.777000+00:00",
                "edited_timestamp": null,
                "flags": 0,
                "components": []
            }
         ]
     }
     ```