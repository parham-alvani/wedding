import type { APIRoute } from "astro";
import { config } from "astro.config";

export const POST: APIRoute = async ({ params, request }) => {
  const resp = await fetch(`${config.backend_url}/guest/${params.id}/answer`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: await request.text(),
  });

  return new Response(resp.body, {
    status: resp.status,
    headers: { "Content-Type": "application/json" },
  });
};
