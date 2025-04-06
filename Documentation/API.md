# API

## Authentication

Each frame has its own public/private keypair.

All calls that indicate Authentication, **must provide the `X-Signature` and `X-Public-Key` HTTP headers**.

`X-Signature` is calculated as:

```
X-SIGNATURE = BASE64( ED25519_SIGN(
  PRIVATE_KEY,
  HTTP_METHOD + "\n" + PATH + "\n" + DATE
))
```

And `X-Public-Key` is:

```
X-PUBLIC-KEY = FRAME_ID:BASE64(public key)
```

If you want to use the admin keypair, and not the frame keypair, set `FRAME_ID` to 0.

The admin key is the key configured in Farma during setup, and it is stored in `config.yaml`

Sample code in Javascript can be found in [examples/nodejs/index.js](../examples/nodejs/index.js)
Check [api/utils.go](../api2/apiClient.go) (used by the Farma CLI tools) for an implementation in Go.

## Endpoints

### Frames

#### Get Frame
|Item|Description |
|:--|:--|
|endpoint| /api/v2/frame/:frameId|
|method | GET|
|authentication| frame or admin |
|GET Parameter| `start`: Used when itterating through paginated results |
|GET Parameter| `limit`: max number of results to fetch |

Returns a JSON object with information about the configured frames. If `id`
is ommited, it will return all frames.

Sample response:

```json
{
  "id": "6m4m2",
  "name": "test-frame",
  "domain": "test.com",
  "publicKey": {
    "frameId": "6m4m2",
    "Key": "s/qu55n1k+sxO5WZ7iHFVapnTVjp0dNRz54jD+pIbhM="
  }
}
```
#### Create Frame
|Item|Description |
|:--|:--|
|endpoint| /api/v2/frame/|
|method | POST|
|authentication| admin |
|payload| `{"name": "frame name", "domain": "frame domain"}`|

It will configure a new frame into Farma.

Sample response:

```json
{
 "frame":{
  "id": "6m4m2",
  "name": "test-frame",
  "domain": "test.com",
  "publicKey": {
    "frameId": "6m4m2",
    "Key": "s/qu55n1k+sxO5WZ7iHFVapnTVjp0dNRz54jD+pIbhM="
  }
},
"private_key": "automatically generated private key",
"public_key": "corresponding public key"
}
```

### Subscriptions

#### Get Subscriptions
|Item|Description |
|:--|:--|
|endpoint| /api/v2/subscription/:frameid|
|method | GET|
|authentication| frame or admin |
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
      "url": "https://api.warpcast.com/v2/frame-notifications",
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
|endpoint| /api/v2/logs/:frameId/:userId|
|method | GET|
|authentication| frame or admin |
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
|endpoint| /api/v2/notification/:frameId|
|method | POST|
|authentication| frame or admin |
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
