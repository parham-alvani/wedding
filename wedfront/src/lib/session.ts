import { createHmac, randomBytes, timingSafeEqual } from "node:crypto";

/**
 * Secret used to sign the "this visitor proved who they are" cookie.
 *
 * Without WEDFRONT_SESSION_SECRET a random one is generated per process: the
 * gate still works, but restarting the server asks everyone for their name
 * again. Set it in production.
 */
const secret = process.env.WEDFRONT_SESSION_SECRET || randomBytes(32).toString("hex");

if (!process.env.WEDFRONT_SESSION_SECRET) {
  console.warn(
    "[wedfront] WEDFRONT_SESSION_SECRET is not set; guest sessions will not survive a restart",
  );
}

/** Name of the cookie holding a verified guest session. */
export const cookieName = "wedfront_guest";

/** A month: long enough to cover the run-up to the wedding. */
export const cookieMaxAge = 60 * 60 * 24 * 30;

function sign(id: string): string {
  return createHmac("sha256", secret).update(id).digest("hex");
}

/** Builds the cookie value proving this visitor opened `id`. */
export function token(id: string): string {
  return `${id}.${sign(id)}`;
}

/** Checks a cookie value against the guest whose page is being rendered. */
export function verifyToken(value: string | undefined, id: string): boolean {
  if (!value) return false;

  const separator = value.lastIndexOf(".");
  if (separator < 1) return false;

  const signed = value.slice(0, separator);
  const digest = value.slice(separator + 1);

  // The signature covers the id, so a cookie earned on one invitation cannot
  // be replayed against another.
  if (signed !== id) return false;

  const expected = Buffer.from(sign(id), "utf8");
  const actual = Buffer.from(digest, "utf8");

  if (expected.length !== actual.length) return false;

  return timingSafeEqual(expected, actual);
}
