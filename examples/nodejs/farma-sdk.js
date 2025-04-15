import * as http from "http";
import * as https from "https";
import { ed25519 } from "@noble/curves/ed25519";
import { Base64 } from "js-base64";

class FarmaSDK {
    constructor(config) {
        this.config = {
            hostname: config.hostname || "127.0.0.1",
            port: config.port || 8080,
            protocol: config.protocol || "http",
            frameId: config.frameId,
            privateKey: config.privateKey,
            // If true, use admin key (frameId = 0)
            isAdmin: config.isAdmin || false
        };

        // Initialize key pair
        // The private key is base64 encoded and contains both private and public parts
        const fullKeyBytes = Base64.toUint8Array(this.config.privateKey);

        if (fullKeyBytes.length < 64) {
            throw new Error("Invalid private key length. Expected at least 64 bytes (32 private + 32 public).");
        }

        // The private key is the first 32 bytes
        this.privateKeyBytes = fullKeyBytes.slice(0, 32);
        // The public key is the last 32 bytes
        this.publicKey = fullKeyBytes.slice(32, 64);
        this.publicKey64 = Base64.fromUint8Array(this.publicKey);
    }

    async request(method, path, body = null) {
        const date = new Date(Date.now()).toGMTString();
        
        // Ensure path starts with /
        let normalizedPath = path;
        if (!normalizedPath.startsWith('/')) {
            normalizedPath = '/' + normalizedPath;
        }

        const options = {
            hostname: this.config.hostname,
            port: this.config.port,
            path: `/api/v2${normalizedPath}`,
            method: method,
            headers: {
                "Content-Type": "application/json",
            },
        };

        // Add body if present
        if (body) {
            options.headers["Content-Length"] = Buffer.byteLength(JSON.stringify(body));
        }

        // Sign request using the private key
        const message = Buffer.from(`${method}\n${options.path}\n${date}`);
        const signature = ed25519.sign(message, this.privateKeyBytes);
        const signature64 = Base64.fromUint8Array(signature);

        // Add authentication headers
        options.headers["X-Signature"] = signature64;
        options.headers["X-Public-Key"] = `${this.config.isAdmin ? "0" : this.config.frameId}:${this.publicKey64}`;
        options.headers["X-Date"] = date;

        console.log(`Making ${method} request to ${this.config.protocol}://${options.hostname}:${options.port}${options.path}`);
        console.log('Headers:', options.headers);
        console.log('X-Public-Key:', options.headers["X-Public-Key"]);

        return new Promise((resolve, reject) => {
            const req = (this.config.protocol === 'https' ? https : http).request(options, (res) => {
                console.log(`Response status: ${res.statusCode} ${res.statusMessage}`);
                console.log('Response headers:', res.headers);

                let data = "";
                res.on("data", (chunk) => {
                    data += chunk;
                });

                res.on("end", () => {
                    try {
                        if (res.statusCode >= 400) {
                            const error = new Error(`HTTP Error ${res.statusCode}: ${res.statusMessage}`);
                            error.statusCode = res.statusCode;
                            error.response = data;
                            reject(error);
                            return;
                        }
                        const response = JSON.parse(data);
                        resolve(response);
                    } catch (e) {
                        const error = new Error(`Failed to parse response: ${e.message}`);
                        error.response = data;
                        reject(error);
                    }
                });
            });

            req.on("error", (err) => {
                console.error('Request error:', err);
                reject(err);
            });

            if (body) {
                req.write(JSON.stringify(body));
            }
            req.end();
        });
    }

    // Frame Methods
    async getFrame(frameId = null) {
        const id = frameId || this.config.frameId;
        if (!id) {
            throw new Error("Frame ID is required");
        }
        return this.request("GET", `frame/${id}`);
    }

    async createFrame(name, domain) {
        if (!this.config.isAdmin) {
            throw new Error("Admin privileges required to create frames");
        }
        return this.request("POST", "frame/", { name, domain });
    }

    async updateFrame(frameId, updates) {
        return this.request("POST", `frame/${frameId}`, updates);
    }

    // Subscription Methods
    async getSubscriptions(frameId = null) {
        const id = frameId || this.config.frameId;
        if (!id) {
            throw new Error("Frame ID is required");
        }
        return this.request("GET", `subscription/${id}`);
    }

    // Notification Methods
    async getNotifications(frameId, notificationId = null) {
        const path = notificationId ? 
            `notification/${frameId}/${notificationId}` : 
            `notification/${frameId}`;
        return this.request("GET", path);
    }

    async sendNotification(frameId, title, body, url = "", userIds = null) {
        const payload = {
            frameId,
            title,
            body,
            url
        };
        if (userIds) {
            payload.userIds = userIds;
        }
        return this.request("POST", `notification/${frameId}`, payload);
    }

    // User Logs Methods
    async getUserLogs(frameId, userId = null) {
        const path = userId ? `logs/${frameId}/${userId}` : `logs/${frameId}`;
        return this.request("GET", path);
    }

    // Database Methods (Admin only)
    async getDbKeys(prefix = "") {
        return this.request("GET", `dbkeys/${prefix}`);
    }
}

export default FarmaSDK; 