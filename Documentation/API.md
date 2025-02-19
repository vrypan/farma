# API

## Request format and authentication

All API calls have the same endpoint, `/api/v1/`, and they follow the
[JSON Farcaster Signatures spec](https://github.com/farcasterxyz/protocol/discussions/208).

However, instead of using a Farcaster AppKey, the requests are signed by a shared key generated
during `fargo setup`. `header.type` is set to `shared` and `fid` is set to `0`.

For example:

`"header":{ "fid": 0, "type": "shared", "key": "0x30a63474db060dc307353490e6417d1318dbb06f3c1c208bbce499b710033aee" }`

In the future, `farma` may support `custody` and `app_key` types, if it makes sense.

## Payload

The actual API call is included in the `payload` field of JFS. It contains two parts, `command` and `params`.

For example:

```
"payload": {
  "command": "frames/get",
  "params": {
    "id": 1
  }
}
```

Once you have the header and the payload, you can construct the JFS object and send it to `/api/v1/`.

Sample Javascript:

```javascript

const nacl = require("tweetnacl");
const util = require("tweetnacl-util");

function toBase64URL(s) {
  if (typeof s === "string") {
    s = new Uint8Array(util.decodeUTF8(s));
  }
  return util
    .encodeBase64(s)
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/, "");
}

const PUB_STR = "0x...." // Private Ed25519 key in hex format
const PRIV_STR = "0x...." // Private Ed25519 key in hex format
const privateKey = new Uint8Array(Buffer.from(PRIV_STR.slice(2), "hex"));
const publicKey = new Uint8Array( Buffer.from(PUB_STR.slice(2), "hex") );

const header = toBase64URL(
  JSON.stringify({ fid: 0, type: "shared", key: PUB_HEX }),
);
const payload = toBase64URL(
  JSON.stringify({
    command: "notification/send",
    params: {
      frame: "farma2",
      title: "Hello there",
      body: "This is the message body",
      url: "",
    },
  }),
);

const S = header + "." + payload; // Construct the message to sign
const message = util.decodeUTF8(S);
const signature = nacl.sign.detached(message, privateKey);
const signatureBase64URL = toBase64URL(signature);

// Construct the JSON
const json = JSON.stringify({
  header: header,
  payload: payload,
  signature: signatureBase64URL,
});

console.log(json);

```

`cmd/cli.go` contains a similar implementation in Go.

## Commands

The following commands are available. This section will definitely change
because they need normalization, pagination, etc.

### `notification/send`
- Expects `frame` (the frame short name), `title`, `body`, `url` in params.
- It will send the notification to all users subscribed to `frame`.

### `frames/get`
- Optional parameter `id` (number, the frameId)
- It will return all frames, or only the specific frame id

### `frames/add`
- Expects `name`, `domain`, `webhook`
- It will configure a new frame.

### `subscriptions/get`
- Optional parameters `frameId` and `limit`.
- It will return subscriptions for all frmaes or only for `frameId`. `limit` is the maximum number of results to return.

### `logs/get`
- Optional parameters `userId` and `limit`.
- It will return activity logs (subscription status updates, notifications sent) for all users or only for `userId`.
`limit` is the maximum number of results to return.
