import express from "express";
import { handler as ssrHandler } from "./dist/server/entry.mjs";
import proxy from "express-http-proxy";
import chalk from "chalk";

const app = express();
const base = "/";
app.use(base, express.static("dist/client/"));
app.use(ssrHandler);
app.use("/api/", proxy("http://127.0.0.1:1378/"));

console.log(chalk.yellowBright.bold("Wedfront, Wedding Frontend Proxy"));
console.log(chalk.greenBright.bold("Listening on :8080"));
console.log("---");
console.log();

app.listen(8080);
