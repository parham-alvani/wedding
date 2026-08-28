import type { APIRoute } from "astro";
import { config } from "astro.config";
import { cookieMaxAge, cookieName, token } from "~/lib/session";

export const POST: APIRoute = async ({ params, request, cookies }) => {
  const id = params.id ?? "";

  const resp = await fetch(`${config.backend_url}/guest/${id}/verify`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: await request.text(),
  });

  if (resp.ok) {
    cookies.set(cookieName, token(id), {
      path: `/guests/${id}`,
      httpOnly: true,
      sameSite: "lax",
      secure: new URL(request.url).protocol === "https:",
      maxAge: cookieMaxAge,
    });
  }

  return new Response(resp.body, {
    status: resp.status,
    headers: { "Content-Type": "application/json" },
  });
};
