import type { EventName, Guest } from "~/types";
import { wedding } from "~/wedding.config";

/**
 * Resolves which ceremonies a guest is invited to. An absent or empty value
 * means all of them, matching the backend, so a guest list that predates
 * tiering keeps working.
 */
export function invitedEvents(guest: Guest): Array<EventName> {
  const all = wedding.events.map((event) => event.key as EventName);

  const raw = (guest.events ?? "").trim();
  if (raw === "") return all;

  const wanted = new Set(
    raw
      .split(/[,|\s]+/)
      .map((part) => part.trim().toLowerCase())
      .filter(Boolean),
  );

  const invited = all.filter((key) => wanted.has(key));

  // An unrecognised value should not leave a guest with no invitation at all.
  return invited.length > 0 ? invited : all;
}

/** The venue and date configuration for each ceremony a guest may attend. */
export function ceremoniesFor(guest: Guest) {
  const invited = new Set(invitedEvents(guest));

  return wedding.events
    .filter((event) => invited.has(event.key as EventName))
    .map((event) => ({
      key: event.key as EventName,
      venue: wedding.venues[event.venue as keyof typeof wedding.venues],
      date: wedding.dates[event.date as keyof typeof wedding.dates] as string,
    }));
}
