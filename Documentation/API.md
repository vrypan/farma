# API

## Authentication

All calls bellow that indicate Authentication, must provide a `X-Signature` HTTP header.

The public key used to sign the request is expected in the `X-Public-Key` HTTP header. This key
must be already configured in Farma. Right now, the API does not check `X-Public-Key`, and
relies in the single key configured during `farma setup`. In the future, when more than one
keys are supported, `X-Public-Key`will be required.

`X-Signature` is calculated as:

```
X-SIGNATURE = ED25519_SIGN(
  PRIVATE_KEY,
  HTTP_METHOD + "\n" + PATH + "\n" + DATE
)
```

Sample code in Javascript can be found in [examples/nodejs/index.js](../examples/nodejs/index.js)
Check [api/utils.go](../api/utils.go) for an implementation in Go.

## Endpoints

### Frames

#### Get Frames
|Item|Description |
|:--|:--|
|endpoint| /api/v1/frames/:id|
|method | GET|
|authentication| required |
|GET Parameter| `start`: Used when itterating through paginated results |
|GET Parameter| `limit`: max number of results to fetch |

Returns a JSON object with information about the configured frames. If `id`
is ommited, it will return all frames.

Sample response:

```json
{
  "result":[
    {
      "id": 1,
      "name": "example",
      "domain": "example.com",
      "webhook": "/f/16c89d6c-4356-4481-accc-18b3a0b49a2b"
    },
    {
      "id": 2,
      "name": "example2",
      "domain": "example2.com",
      "webhook": "/f/8dd2b175-83ba-48ba-9dd0-41453b4f86ef"
    }
  ],
  "next": "bDp1c2VyOjI4MDoyOjE3NDAyNTU5NzgK"
}
```
#### Create Frame
|Item|Description |
|:--|:--|
|endpoint| /api/v1/frames/|
|method | POST|
|authentication| required |
|payload| `{"name": "frame name", "domain": "frame domain", "webhook":""}`|

It will configure a new frame into Farma. It is advised to leave `webhook` empty, unless you know exactly what you are doing.

Sample response:

```json
  {
    "id": 1,
    "name": "example",
    "domain": "example.com",
    "webhook": "/f/16c89d6c-4356-4481-accc-18b3a0b49a2b"
  },
```

### Subscriptions

#### Get Subscriptions
|Item|Description |
|:--|:--|
|endpoint| /api/v1/subscriptions/:frameid|
|method | GET|
|authentication| required |
|GET Parameter| `start`: Used when itterating through paginated results |
|GET Parameter| `limit`: max number of results to fetch |

This endpoint will return a list of subscriptions. If `frameId` is provided
it will return only subscriptions for that frame.

Sample response:

```json
{
  "result":[
    {
      "frameId": 1,
      "userId": 20396,
      "appId": 9152,
      "status": 2,
      "url": "https://api.warpcast.com/v1/frame-notifications",
      "token": "01952cfc-cc2f-8b3c-4ed8-a476d9b05050",
      "ctime": {
        "seconds": 1739953526,
        "nanos": 954099000
      },
      "mtime": {
        "seconds": 1740216525,
        "nanos": 388740000
      },
      "verified": true,
      "appKey": "aaNtJNyy5/aE0aWssFSjq1EuP6ZU9bHcE53LLxmAEM0="
    },
  ...
  ],
  "next": "bDp1c2VyOjI4MDoyOjE3NDAyNTU5NzgK"
}
```

### User Logs

#### Get Logs
|Item|Description |
|:--|:--|
|endpoint| /api/v1/logs/:userId|
|method | GET|
|authentication| required |
|GET Parameter| `start`: Used when itterating through paginated results |
|GET Parameter| `limit`: max number of results to fetch |

This endpoint will return history logs. If `userId` (fid) is provided
it will return only logs for that user. Events include frame add/remove,
notifications enabled/disabled, and notifications sent.

|EvtType|Description |
|:--|:--|
|1|Frame added|
|2|Frame removed|
|3|Notification enabled|
|4|Notification disabled|
|5|Notification sent|
|6|Notification failed other|
|7|Notification failed invalid|
|8|Notification failed rate limit|

Sample response:

```json
{
  "result": [
    {
      "frameId": 1,
      "userId": 280,
      "appId": 9152,
      "evtType": 2,
      "ctime": {
        "seconds": 1739953487,
        "nanos": 955840000
      }
    },
    {
      "frameId": 1,
      "userId": 280,
      "appId": 9152,
      "evtType": 1,
      "ctime": {
        "seconds": 1739953499,
        "nanos": 487717000
      }
    },
  ...
  ],
  "next": "bDp1c2VyOjI4MDoyOjE3NDAyNTU5NzgK"
}
```

### Notifications

#### Send notification
|Item|Description |
|:--|:--|
|endpoint| /api/v1/notifications/|
|method | POST|
|authentication| required |
|payload| `{"frameId": frameId, "title": "notif title", "body": "notif body", "url": "notification link"}`|

It will send a notification to all subscriebrs of frame `frameId`.
- `url` must be under the frame domain
- You can leave `url` empty (""), to link to the frame itself.
- All users will get the same notification.

It will rerturn the `notificationId` used to send the noitification and the number of notifications
**attempted** to send.

Sample response:

```json
  {
    "NotificationId": "12345678-4356-4481-accc-18b3a0b49a2b"
    "Count" : 3
  },
```
