const { execSync } = require("child_process");

execSync("cd .. && npm run build", { stdio: "inherit" });
