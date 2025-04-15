import FarmaSDK from "./farma-sdk.js";

// Initialize the SDK with a valid private key
// This is a sample key - replace with your actual key
const farma = new FarmaSDK({
  hostname: "127.0.0.1",
  port: 8080,
  protocol: "http", // or "https" for secure connections
  frameId: "z0001", // Updated to match test frame ID
  // This is a sample 64-byte private key (32 bytes private + 32 bytes public)
  privateKey:
    "bSyiLSZtM7/WPUOfgmboyaQsgVhJthYDm2PQFOrhm2UMC8KVJAx2udltMK3M020HMdZO69OUcJIBpaX1ylR/fg==",
});

// Example: Get frame information
async function getFrameInfo() {
  try {
    console.log("\n=== Getting frame info ===");
    const response = await farma.getFrame();
    console.log("Frame Info:", JSON.stringify(response, null, 2));
  } catch (error) {
    console.error("Error getting frame:", error.message);
    if (error.response) {
      console.error("Response:", error.response);
    }
    if (error.statusCode) {
      console.error("Status Code:", error.statusCode);
    }
  }
}

// Example: Send a notification
async function sendNotification() {
  try {
    console.log("\n=== Sending notification ===");
    const response = await farma.sendNotification(
      "z0001", // frameId
      "Test Notification", // title
      "This is a test message", // body
      "https://example.com", // url (optional)
      [123, 456], // userIds (optional)
    );
    console.log("Notification sent:", JSON.stringify(response, null, 2));
  } catch (error) {
    console.error("Error sending notification:", error.message);
    if (error.response) {
      console.error("Response:", error.response);
    }
    if (error.statusCode) {
      console.error("Status Code:", error.statusCode);
    }
  }
}

// Example: Get user logs
async function getUserLogs() {
  try {
    console.log("\n=== Getting user logs ===");
    const response = await farma.getUserLogs("z0001", 123);
    console.log("User Logs:", JSON.stringify(response, null, 2));
  } catch (error) {
    console.error("Error getting user logs:", error.message);
    if (error.response) {
      console.error("Response:", error.response);
    }
    if (error.statusCode) {
      console.error("Status Code:", error.statusCode);
    }
  }
}

// Example: Get notifications
async function getNotifications() {
  try {
    console.log("\n=== Getting notifications ===");
    const response = await farma.getNotifications("z0001");
    console.log("Notifications:", JSON.stringify(response, null, 2));
  } catch (error) {
    console.error("Error getting notifications:", error.message);
    if (error.response) {
      console.error("Response:", error.response);
    }
    if (error.statusCode) {
      console.error("Status Code:", error.statusCode);
    }
  }
}

// Example: Admin operations
async function adminOperations() {
  const adminFarma = new FarmaSDK({
    hostname: "127.0.0.1",
    port: 8080,
    protocol: "http", // or "https" for secure connections
    // This is a sample admin key - replace with your actual admin key
    privateKey:
      "bSyiLSZtM7/WPUOfgmboyaQsgVhJthYDm2PQFOrhm2UMC8KVJAx2udltMK3M020HMdZO69OUcJIBpaX1ylR/fg==",
    isAdmin: true,
  });

  try {
    console.log("\n=== Creating new frame ===");
    const newFrame = await adminFarma.createFrame("New Frame", "example.com");
    console.log("New Frame:", JSON.stringify(newFrame, null, 2));

    console.log("\n=== Getting database keys ===");
    const dbKeys = await adminFarma.getDbKeys();
    console.log("Database Keys:", JSON.stringify(dbKeys, null, 2));
  } catch (error) {
    console.error("Error in admin operations:", error.message);
    if (error.response) {
      console.error("Response:", error.response);
    }
    if (error.statusCode) {
      console.error("Status Code:", error.statusCode);
    }
  }
}

// Run all examples
async function runExamples() {
  console.log("Starting Farma SDK examples...");
  console.log("Using frame ID:", farma.config.frameId);
  console.log("Using protocol:", farma.config.protocol);
  
  await getFrameInfo();
  await sendNotification();
  await getUserLogs();
  await getNotifications();
  await adminOperations();
  
  console.log("\nAll examples completed!");
}

runExamples().catch((error) => {
  console.error("Fatal error:", error);
  process.exit(1);
});
