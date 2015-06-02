# Push API

## App

### POST /apps

App Model:

        {
            "id": "app id",
            "platforms": {
                "apns": {
                    "certificate": "----PEM encoded certificate----",
                    "privateKey": "apns private key"
                },
                "gcm": {
                    "projectId": "...",
                    "apiKey": "...."
                }
            }
        }


### PUT /apps/{appId}

Update app. Request body is the same as App Model.

### GET /apps/{appId}

...


## Subscribers

Subscriber Model:

    {
        "id": "Subscriber Id",
        //? "channels": [],
        "createdAt": "unix timestamp",
        "devices": [
            {
                "platform": "gcm or apns",
                "token": "token is deviceToken in apns, registrationId in gcm",
                "createdAt": "unix timestamp"
            },
            ...
        ]
    }

Subscriber model serves two purposes:

* Abstracts the mapping between a user and her devices. As a results clients do not need to know about device tokens etc to send push notifications.
* Groups multiple devices of a user

### POST /apps/{appId}/devices

**FIX**: ```/apps/{appId}/subscribers/{subscriberId}/devices``` is semantically better but syntactically worse IMHO. 

    {
        "subscriberId": "client defined subscriber Id.",
        "platform": "gcm or apns",
        "token": "token is deviceToken in apns, registrationtId in gcm"
    }

## Channels

### POST /apps/{appId}/channels

Create channel

    {
        "id": "channel id"
    }

**FIX**: What about creating channels on the fly while adding subscribers?


### DELETE /apps/{appId}/channels/{channelId}

Delete a channel


### POST /apps/{appId}/channels/{channelId}/subscribers

Add subscribers to channel:

    {
        "subscribers": ["list", "of", "subscriber", "ids"]
    }


### DELETE /apps/{appId}/channels/{channelId}/subscribers/{subscriberId}

Remove a subscriber from a channel.

**FIX**: Do we really need this?



## Publish

### POST /apps/{appId}/publish

**FIX**: ```publish``` is not a name. Other alternative is push which is also a verb. Noun alternative is ```notifications``` but it is not as good(descriptive) as the former ones imho.

Request:

        {
            "recipients": ["list", "of", "subscriber", "ids", "..."],
            "channels": ["list", "of", "channels"],
            "message": {
                "gcm": {
                    // message for gcm
                },
                "apns": {
                    // message for apn
                }
            }
        }

Response:

    {
        "transactionId": "transaction uuid"
    }


#### Flow

Master:

1. Expand recipients and channels to a list of deviceTokens
2. Push each **(transactionId, platform, deviceToken, message)** quadruplets to message queue.
3. Return transactionId, # of messages to deliver.

Worker:

1. Pull messages
2. Send to appropriate push backend
3. Unregister & log failing subscribers
4. Update transaction counters

