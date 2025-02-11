**Don't use it yet, it's still pre-alpha.**

![farma-git-social](farma-git-social.png)

# farma

**farma** (**FA**rcaster **R**elationships **MA**nagement) is built to make Farcaster Frames v2
notifications easy, and help you get insights on how users interact with your
frames notifications.

- When did a user first add your frame?
- How many notifications have you sent them since?
- When did they enable/disable notifications
- Which Farcaster clients are they using to interact with your frame?

# TODO

- [x] Support multiple frames and multiple clients per user
- [x] Handle subscriptions/unsubscriptions
- [x] Group FIDs in batches of 100 when sending notifications
- [x] Handle Warpcast responses when sending notifications
- [x] Validate subscription signatures
- [x] Validate AppKeys used to sign subscription signatures
- [x] Log user subscription/unsubscriptions
- [ ] Installation instructions
- [ ] Support access to nodes that require authentication
- [ ] Guide admins on how to use the commands, add checks
- [ ] Expose REST API
- [ ] Webhooks?
- [ ] Log notifications sent

# Setup

(This will improve a lot, but in case someone wants to give it an early try.)

1. Download the binary corresponding to your system from `https://github.com/vrypan/farma/releases`.
If you're using macOS, use `brew install vrypan/farma/farma`

2. Run `farma setup`.

3. Create a frame callbackUrl: `farma frame add <myframe> --url=<frame url>`.
`<myframe>` is a short name used to identify your frame. The optional `<frame url>` is the URL of your frame.

4. You will get a relative endpoint, something like `/f/a2b01541-778d-4a2b-9375-8232c70a6ddf`.
You can use `farma frame ls` to see all your endpoints.

5. Add the endpoint (including your server name) in the frame's `.well-known/farcaster.json` callbackUrl.

6. Start the farma server: `farma server`.

7. You can send a notification to all users subscribed to a frame, using `farma notify <frame name>`.
