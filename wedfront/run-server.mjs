import express from 'express';
import { handler as ssrHandler } from './dist/server/entry.mjs';

const app = express();
const base = '/';
app.use(base, express.static('dist/client/'));
app.use(ssrHandler);

app.listen(8080);