import { wedding } from "~/wedding.config";

/** Formats a date as an iCalendar UTC timestamp: 20240616T150000Z. */
function stamp(date: Date): string {
  return `${date.toISOString().replace(/[-:]/g, "").split(".")[0]}Z`;
}

/**
 * iCalendar forbids raw commas, semicolons and newlines in text values, and
 * folds lines longer than 75 octets.
 */
function escapeText(value: string): string {
  return value
    .replace(/\\/g, "\\\\")
    .replace(/([,;])/g, "\\$1")
    .replace(/\n/g, "\\n");
}

function fold(line: string): string {
  if (line.length <= 75) return line;

  const chunks: Array<string> = [line.slice(0, 75)];
  for (let i = 75; i < line.length; i += 74) {
    chunks.push(` ${line.slice(i, i + 74)}`);
  }

  return chunks.join("\r\n");
}

export interface Event {
  /** Stable identifier, so re-importing updates rather than duplicates. */
  uid: string;
  start: Date;
  title: string;
  location: string;
  description: string;
}

/** Renders one or more events as an iCalendar document. */
export function ics(events: Array<Event>): string {
  const lines: Array<string> = [
    "BEGIN:VCALENDAR",
    "VERSION:2.0",
    "PRODID:-//wedding//EN",
    "CALSCALE:GREGORIAN",
    "METHOD:PUBLISH",
  ];

  for (const event of events) {
    const end = new Date(event.start.getTime() + wedding.calendar.durationHours * 60 * 60 * 1000);

    lines.push(
      "BEGIN:VEVENT",
      `UID:${event.uid}`,
      // A fixed DTSTAMP keeps the file byte-identical between requests.
      `DTSTAMP:${stamp(event.start)}`,
      `DTSTART:${stamp(event.start)}`,
      `DTEND:${stamp(end)}`,
      `SUMMARY:${escapeText(event.title)}`,
      `LOCATION:${escapeText(event.location)}`,
      `DESCRIPTION:${escapeText(event.description)}`,
      "END:VEVENT",
    );
  }

  lines.push("END:VCALENDAR");

  return `${lines.map(fold).join("\r\n")}\r\n`;
}

/** The two ceremonies, straight from the wedding config. */
export function ceremonies(): Array<Event> {
  const { couple, dates, venues, site } = wedding;
  const names = `${couple.wife.name} & ${couple.husband.name}`;

  return [
    {
      uid: `wedding@${new URL(site.url).host}`,
      start: new Date(dates.wedding),
      title: `${names} — Wedding`,
      location: `${venues.wedding.name}, ${venues.wedding.address}`,
      description: venues.wedding.whenEnglish || `${names} wedding`,
    },
  ];
}
