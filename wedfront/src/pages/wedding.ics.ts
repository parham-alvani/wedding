import type { APIRoute } from "astro";
import { ceremonies, ics } from "~/lib/calendar";

export const GET: APIRoute = () =>
  new Response(ics(ceremonies()), {
    headers: {
      "Content-Type": "text/calendar; charset=utf-8",
      "Content-Disposition": 'attachment; filename="wedding.ics"',
    },
  });
