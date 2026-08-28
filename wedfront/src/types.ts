export interface FooterLink {
  description: string;
  icon: string;
  url: string;
}

export interface Answer {
  coming: boolean;
  plus_one: boolean;
  dietary?: string;
  song?: string;
}

/** A ceremony key, matching the backend's model.Event. */
export type EventName = "engagement" | "wedding";

export interface Guest {
  first_name: string;
  last_name: string;
  id: string;
  is_family?: boolean;
  spouse_first_name?: string;
  spouse_last_name?: string;
  /** Comma-separated ceremonies; absent or empty means all of them. */
  events?: string;
  answer?: Answer;
}
