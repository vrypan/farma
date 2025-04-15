# Farma Node.js SDK

A Node.js SDK for interacting with the Farma API.

## Installation

```bash
npm install @noble/curves js-base64
```

## Usage

### Basic Setup

```javascript
import FarmaSDK from './farma-sdk.js';

const farma = new FarmaSDK({
    hostname: "127.0.0.1",  // Farma server hostname
    port: 8080,            // Farma server port
    frameId: "z0001",      // Your frame ID
    privateKey: "your-private-key-here"  // Base64 encoded private key
});
```

### Frame Operations

```javascript
// Get all frames
const frames = await farma.getFrame();

// Get specific frame
const frame = await farma.getFrame("z0001");

// Create a new frame (admin only)
const adminFarma = new FarmaSDK({
    hostname: "127.0.0.1",
    port: 8080,
    privateKey: "admin-private-key-here",
    isAdmin: true
});
const newFrame = await adminFarma.createFrame("New Frame", "example.com");

// Update frame
const updatedFrame = await farma.updateFrame("z0001", {
    name: "Updated Name",
    domain: "updated.com"
});
```

### Notification Operations

```javascript
// Send a notification
const notification = await farma.sendNotification(
    "z0001",                    // frameId
    "Test Notification",        // title
    "This is a test message",   // body
    "https://example.com",      // url (optional)
    [123, 456]                  // userIds (optional)
);

// Get notifications
const notifications = await farma.getNotifications("z0001");

// Get specific notification
const notification = await farma.getNotifications("z0001", "notification-id");
```

### Subscription Operations

```javascript
// Get all subscriptions
const subscriptions = await farma.getSubscriptions();

// Get frame subscriptions
const frameSubscriptions = await farma.getSubscriptions("z0001");
```

### User Logs

```javascript
// Get all user logs for a frame
const logs = await farma.getUserLogs("z0001");

// Get logs for specific user
const userLogs = await farma.getUserLogs("z0001", 123);
```

### Admin Operations

```javascript
const adminFarma = new FarmaSDK({
    hostname: "127.0.0.1",
    port: 8080,
    privateKey: "admin-private-key-here",
    isAdmin: true
});

// Get database keys
const dbKeys = await adminFarma.getDbKeys();
```

## Error Handling

All SDK methods return Promises and can be used with try/catch:

```javascript
try {
    const response = await farma.getFrame();
    console.log(response);
} catch (error) {
    console.error("Error:", error);
}
```

## Authentication

The SDK handles authentication automatically using Ed25519 signatures. You need to provide:

1. A private key (base64 encoded)
2. Frame ID (unless using admin key)
3. Admin flag (if using admin key)

The SDK will:
- Generate the public key from the private key
- Sign requests with the current timestamp
- Add all required headers (X-Signature, X-Public-Key, X-Date)

## Dependencies

- `@noble/curves`: For Ed25519 signature generation
- `js-base64`: For base64 encoding/decoding

## Example

See `example.js` for a complete example of using the SDK. 