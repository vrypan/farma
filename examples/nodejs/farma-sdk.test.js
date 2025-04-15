import { Base64 } from "js-base64";
import { ed25519 } from "@noble/curves/ed25519";
import FarmaSDK from "./farma-sdk.js";

// Test data
const TEST_PRIVATE_KEY = "qFew+F/eBNwXOZPiHoiWsFI0B7pTnXBvXTeklCnu8GMaT4EXouj3p1T6PYw8wJzp8Bp2qovZ1wMiH86ABk8qBQ==";
const TEST_FRAME_ID = "z0002";
const EXPECTED_PUBLIC_KEY = "Gk+BF6Lo96dU+j2MPMCc6fAadqqL2dcDIh/OgAZPKgU=";

// Test key handling
console.log("Testing key handling...");

// 1. Test full key decoding
const fullKeyBytes = Base64.toUint8Array(TEST_PRIVATE_KEY);
console.log("Full key length:", fullKeyBytes.length);
console.log("Full key (hex):", Buffer.from(fullKeyBytes).toString('hex'));

// 2. Test private key extraction
const privateKeyBytes = fullKeyBytes.slice(0, 32);
console.log("Private key length:", privateKeyBytes.length);
console.log("Private key (hex):", Buffer.from(privateKeyBytes).toString('hex'));

// 3. Test public key generation
const generatedPublicKey = ed25519.getPublicKey(privateKeyBytes);
const generatedPublicKey64 = Base64.fromUint8Array(generatedPublicKey);
console.log("Generated public key (base64):", generatedPublicKey64);
console.log("Expected public key (base64):", EXPECTED_PUBLIC_KEY);

// 4. Test SDK initialization
const sdk = new FarmaSDK({
    frameId: TEST_FRAME_ID,
    privateKey: TEST_PRIVATE_KEY
});

// 5. Test request signing
const date = new Date(Date.now()).toGMTString();
const method = "GET";
const path = `/api/v2/frame/${TEST_FRAME_ID}`;
const message = `${method}\n${path}\n${date}`;
console.log("\nSignature calculation:");
console.log("Method:", method);
console.log("Path:", path);
console.log("Date:", date);
console.log("Message to sign:", message);
console.log("Message bytes (hex):", Buffer.from(message).toString('hex'));

const signature = ed25519.sign(Buffer.from(message), privateKeyBytes);
console.log("Signature bytes (hex):", Buffer.from(signature).toString('hex'));
const signature64 = Base64.fromUint8Array(signature);
console.log("Signature (base64):", signature64);

// 6. Test signature verification
const isValid = ed25519.verify(signature, Buffer.from(message), generatedPublicKey);
console.log("Signature verification:", isValid);

// 7. Test full X-Public-Key header
const expectedXPublicKey = `${TEST_FRAME_ID}:${EXPECTED_PUBLIC_KEY}`;
const generatedXPublicKey = `${TEST_FRAME_ID}:${generatedPublicKey64}`;
console.log("\nX-Public-Key header:");
console.log("Expected:", expectedXPublicKey);
console.log("Generated:", generatedXPublicKey);

// 8. Verify results
const publicKeyMatch = generatedPublicKey64 === EXPECTED_PUBLIC_KEY;
console.log("\nTest Results:");
console.log("Public key matches expected:", publicKeyMatch);
if (!publicKeyMatch) {
    console.log("ERROR: Public key does not match expected value!");
    console.log("Expected:", EXPECTED_PUBLIC_KEY);
    console.log("Got:     ", generatedPublicKey64);
}

// 9. Test SDK request
console.log("\nTesting SDK request...");
sdk.request("GET", `frame/${TEST_FRAME_ID}`)
    .then(response => {
        console.log("Request successful:", response);
    })
    .catch(error => {
        console.error("Request failed:", error);
    }); 