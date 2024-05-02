import express from "express";
import { handler as ssrHandler } from "./dist/server/entry.mjs";
import proxy from "express-http-proxy";

const app = express();
const base = "/";
app.use(base, express.static("dist/client/"));
app.use(ssrHandler);
app.use("/api/", proxy("http://127.0.0.1:1378/"));

app.listen(8080);
