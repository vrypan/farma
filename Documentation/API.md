# API

A nodejs SDK is available in [examples/nodejs](../examples/nodejs).

If you are using Go, check `apiv2/apiClient.go`

## Authentication

Each frame has its own public/private keypair. There are two types of authentication:

1. **Admin Authentication**: Uses the admin keypair configured in Farma during setup (stored in `config.yaml`)
2. **Frame Authentication**: Uses the frame's own keypair

All authenticated requests must provide the following HTTP headers:

| Header | Description |
|--------|-------------|
| `X-Signature` | Base64 encoded Ed25519 signature of `METHOD\nPATH\nDATE` |
| `X-Public-Key` | Format: `FRAME_ID:BASE64(public_key)`. Use `0` as FRAME_ID for admin key |
| `X-Date` | Current date in RFC1123 format (must be within 10 seconds of server time) |

The signature is verified using Ed25519. For frame authentication, the public key must exist in the database.

## Access Control

Each endpoint has an Access Control Level (ACL):

1. `ACL_ADMIN`: Only admin key is allowed
2. `ACL_FRAME_OR_ADMIN`: Both frame key and admin key are allowed

When using frame authentication, the `frameId` in the request must match the frame ID in the public key.

## Endpoints

### Frames

#### Get Frame
|Item|Description |
|:--|:--|
|endpoint| `/api/v2/frame/:frameId`|
|method | GET|
|authentication| `ACL_FRAME_OR_ADMIN` |
|GET Parameter| `start`: Used when iterating through paginated results |
|GET Parameter| `limit`: max number of results to fetch |

Returns a JSON object with information about the configured frames. If `frameId` is omitted, it will return all frames.

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
|endpoint| `/api/v2/frame/`|
|method | POST|
|authentication| `ACL_ADMIN` |
|payload| `{"name": "frame name", "domain": "frame domain"}`|

Creates a new frame and generates a new keypair.

Sample response:

```json
{
  "frame": {
    "id": "6m4m2",
    "name": "test-frame",
    "domain": "test.com",
    "publicKey": {
      "frameId": "6m4m2",
      "Key": "s/qu55n1k+sxO5WZ7iHFVapnTVjp0dNRz54jD+pIbhM="
    }
  },
  "private_key": "automatically generated private key",
  "public_key": "X-Public-Key header to be used in future requests"
}
```

#### Update Frame
|Item|Description |
|:--|:--|
|endpoint| `/api/v2/frame/:frameId`|
|method | POST|
|authentication| `ACL_FRAME_OR_ADMIN` |
|payload| `{"name": "new name", "domain": "new domain", "public_key": "new public key"}`|

Updates an existing frame's properties. All fields are optional.

### Subscriptions

#### Get Subscriptions
|Item|Description |
|:--|:--|
|endpoint| `/api/v2/subscription/:frameId`|
|method | GET|
|authentication| `ACL_FRAME_OR_ADMIN` |
|GET Parameter| `start`: Used when iterating through paginated results |
|GET Parameter| `limit`: max number of results to fetch |

Returns a list of subscriptions. If `frameId` is provided, it returns only subscriptions for that frame.

Sample response:

```json
{
  "result": [
    {
      "frameId": "6m4m2",
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
    }
  ],
  "next": "bDp1c2VyOjI4MDoyOjE3NDAyNTU5NzgK"
}
```

#### Import Subscriptions
|Item|Description |
|:--|:--|
|endpoint| `/api/v2/subscription-import/:frameId`|
|method | POST|
|authentication| `ACL_FRAME_OR_ADMIN` |
|payload| |
|POST Parameter| `appId`: FID of the application. Default:9152 |
|POST Parameter| `appUrl`: Notifications endpoint for appId. Default: https://api.warpcast.com/v1/frame-notifications |

Imports a subscriptions from a CSV file. The first two columns must be `fid` and `subscriptionToken`.

Returns the number of subscriptions imported

Sample response:

```json
{
  "result":
    {
      "entries": 100,
      "message":"Subscriptions imported successfully",
    }
}
```


### Notifications

#### Send Notification
|Item|Description |
|:--|:--|
|endpoint| `/api/v2/notification/:frameId`|
|method | POST|
|authentication| `ACL_FRAME_OR_ADMIN` |
|payload| `{"frameId": "frameId", "title": "notif title", "body": "notif body", "url": "notification link", "userIds": [123, 456]}`|

Sends a notification to subscribers of frame `frameId`:
- If `userIds` is provided, sends only to those users
- Otherwise sends to all subscribers
- `url` must be under the frame domain
- If `url` is empty, links to the frame itself

Sample response:

```json
{
  "NotificationId": "12345678-4356-4481-accc-18b3a0b49a2b",
  "NotificationVersion": 1,
  "Count": 3
}
```

#### Get Notifications
|Item|Description |
|:--|:--|
|endpoint| `/api/v2/notification/:frameId/:notificationId`|
|method | GET|
|authentication| `ACL_FRAME_OR_ADMIN` |
|GET Parameter| `start`: Used when iterating through paginated results |
|GET Parameter| `limit`: max number of results to fetch |

Returns notifications for a frame, optionally filtered by notificationId. Each notification includes its version history.

Sample response:

```json
{
  "result": [
    {
      "frameId": "6m4m2",
      "appId": 9152,
      "id": "12345678-4356-4481-accc-18b3a0b49a2b",
      "endpoint": "https://api.warpcast.com/v2/frame-notifications",
      "title": "Test Notification",
      "message": "This is a test notification",
      "link": "https://test.com/frame",
      "tokens": {
        "token1": 123,
        "token2": 456
      },
      "successTokens": ["token1"],
      "failedTokens": ["token2"],
      "rateLimitedTokens": [],
      "version": 1,
      "ctime": {
        "seconds": 1739953487,
        "nanos": 955840000
      }
    }
  ],
  "next": "bDp1c2VyOjI4MDoyOjE3NDAyNTU5NzgK"
}
```

### User Logs

#### Get Logs
|Item|Description |
|:--|:--|
|endpoint| `/api/v2/logs/:frameId/:userId`|
|method | GET|
|authentication| `ACL_FRAME_OR_ADMIN` |
|GET Parameter| `start`: Used when iterating through paginated results |
|GET Parameter| `limit`: max number of results to fetch |

Returns history logs. If `userId` is provided, returns only logs for that user.

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
      "frameId": "6m4m2",
      "userId": 280,
      "appId": 9152,
      "evtType": 2,
      "ctime": {
        "seconds": 1739953487,
        "nanos": 955840000
      }
    }
  ],
  "next": "bDp1c2VyOjI4MDoyOjE3NDAyNTU5NzgK"
}
```

### Database

#### Get Keys
|Item|Description |
|:--|:--|
|endpoint| `/api/v2/dbkeys/:prefix`|
|method | GET|
|authentication| `ACL_ADMIN` |
|GET Parameter| `start`: Used when iterating through paginated results |
|GET Parameter| `limit`: max number of results to fetch |

Returns database keys matching the prefix.

### System

#### Version
|Item|Description |
|:--|:--|
|endpoint| `/api/v2/version`|
|method | GET|
|authentication| None |
|response| `{"version": "1.0.0"}` |

Returns the current Farma version.

#### New Keypair
|Item|Description |
|:--|:--|
|endpoint| `/api/v2/keypair/:frameId`|
|method | GET|
|authentication| `ACL_FRAME_OR_ADMIN` |
|response| `{"private_key": "base64", "public_key": "base64"}` |

Utility endpoint that generates a new keypair. It does not affect the database in any way.
