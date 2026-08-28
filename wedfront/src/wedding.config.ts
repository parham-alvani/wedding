/**
 * Every guest-facing string, date, link and asset name lives here.
 *
 * Forking this project for your own wedding should mean editing this file,
 * replacing `src/assets/wedding.jpg` with your own photo, dropping your two
 * tracks into `public/`, and nothing else.
 */
export const wedding = {
  couple: {
    husband: {
      name: "Parham",
      /** Name in the ceremony locale (see `ceremony` below). */
      nameLocal: "پرهام",
      lastName: "Alvani",
      emoji: "🐼",
      socials: {
        github: "https://github.com/parham-alvani",
        instagram: "https://instagram.com/1995parham",
      },
    },
    wife: {
      name: "Elaheh",
      nameLocal: "الهه",
      lastName: "Dastan",
      emoji: "🥥",
      socials: {
        github: "https://github.com/elahe-dastan",
        instagram: "https://instagram.com/elahe.dstn",
      },
    },
    lastNameLocal: "داستان - الوانی",
  },

  dates: {
    wedding: "Jun 16, 2024 18:30:00+03:30",
    engaged: "May 10, 2024 19:00:00+03:00",
    /**
     * After this moment the RSVP form is replaced by a closed notice. Leave
     * empty to keep it open forever. Keep it in step with
     * `wedding.rsvp_deadline` in the backend, which does the enforcing.
     */
    rsvpDeadline: "",
  },

  /** How long each ceremony runs, for the "add to calendar" links. */
  calendar: {
    durationHours: 5,
  },

  site: {
    url: "https://parham-alvani.github.io/wedding",
    github: "https://github.com/parham-alvani/wedding",
  },

  music: {
    ceremony: "For-the-Rest-of-My-Life.mp3",
    engaged: "Too-Soon.mp3",
  },

  /**
   * Locale of the two ceremony pages (`/wedding` and `/engaged`). Drives the
   * `lang`/`dir` attributes, the invitation font, and the countdown numerals.
   * The landing page and the per-guest pages are always English.
   */
  ceremony: {
    lang: "fa",
    dir: "rtl",
  },

  /** Invitation card text, in the ceremony locale. */
  invitations: {
    wedding: {
      lines: ["به نام خالق عشق", "عشق بهانه آغاز بود، بهانه سبز باهم زیستن"],
      /** Rendered under the lines; omit to hide. */
      signature: "داستان - الوانی",
    },
    engagement: {
      lines: ["به سنت عشق گرد هم می‌آییم", "آنجا که دوست داشتن تنها کلام زندگی است"],
      signature: "",
    },
  },

  /** The English letter and RSVP form shown on each guest's page. */
  guestLetter: {
    /**
     * Ask visitors to confirm whose invitation this is before showing it.
     * The page is rendered on the server, so nothing personal reaches the
     * browser until the name matches. Either partner's first or last name is
     * accepted. This is a soft gate against a forwarded link, not a password.
     */
    verify: {
      enabled: true,
      heading: "Who is this?",
      prompt: "Type your first or last name to open your invitation.",
      placeholder: "Your name",
      submit: "Open invitation",
      wrong: "That name does not match this invitation.",
      tooMany: "Too many tries. Please wait a few minutes and try again.",
      failed: "Something went wrong. Please try again.",
    },
    body: "We've planned the most special day of our life and we'd be thrilled to see you in our ceremony where we want to celebrate love",
    signOff: "Best,",
    /** Shown instead of the form once the deadline has passed. */
    rsvpClosed: "The RSVP has closed. Please call us and we will sort it out.",
    /** Shown under the form while guests can still change their answer. */
    changeHint: "Changed your mind? You can update this until the RSVP closes.",
    addToCalendar: "Add to calendar",
    rsvp: {
      heading: "RSVP",
      accept: "Accept with pleasure",
      decline: "Decline with regret",
      plusOne: "Plus 1 adult",
      submit: "Respond to invitation",
      update: "Update your answer",
      submitted: "Already responded",
    },

    /**
     * Optional free-text fields on the RSVP. Set `enabled: false` to drop
     * either one; guests who are marked as family never see the coming /
     * declining question but still get these, because the kitchen and the DJ
     * need them just the same.
     */
    extras: {
      dietary: {
        enabled: true,
        label: "Anything we should tell the kitchen?",
        placeholder: "Allergies, vegetarian, …",
      },
      song: {
        enabled: true,
        label: "A song you would love to hear",
        placeholder: "Artist — Title",
      },
      /** Shown to family guests, who have nothing to RSVP to. */
      familyHeading: "A couple of things before the night",
      familySubmit: "Save",
      familySaved: "Saved — thank you!",
    },
  },

  /**
   * The ceremonies, in the order they happen. The keys match the backend's
   * event names, and a guest's page shows only the ones they are invited to.
   */
  events: [
    { key: "engagement", venue: "engagement", date: "engaged" },
    { key: "wedding", venue: "wedding", date: "wedding" },
  ],

  venues: {
    wedding: {
      name: "Baran Garden Hall",
      /** When and where, in the ceremony locale. */
      when: "یکشنبه ۲۷ خرداد از ساعت ۱۸:۳۰ عمارت باران",
      address: "گرمدره، خیابان سعادتیه، کوچه تیر، پلاک ۲۲",
      /** Same thing in English, for the guest pages. */
      whenEnglish: "Hope to see you on Khordad 27th st 6:30 p.m. at Baran mansion garden",
      maps: {
        google: "https://maps.app.goo.gl/tPnYeA5QyvA9LYhEA",
        /** Iranian maps service. Leave empty to hide the link. */
        neshan: "https://nshn.ir/32_bvZnM5x5OAH",
        embed:
          "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3238.355555313541!2d51.07962767712217!3d35.742064672567736!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x3f8dec16562faf13%3A0xb269091e86076da0!2sBaran%20Garden%20Hall!5e0!3m2!1sen!2sat!4v1714895716360!5m2!1sen!2sat",
      },
    },
    engagement: {
      name: "Noghre Hall",
      when: "جمعه ۲۱ام اردیبهشت ماه از ساعت ۱۹",
      address: "تهران، اندرزگو، بلوار صبا، نبش خیابان کریمی، پلاک ۳، واحد ۳",
      whenEnglish: "Hope to see you on Ordibehesht 21st at 7 p.m. at Noghre hall",
      maps: {
        google: "https://maps.app.goo.gl/4r9Rj1sAEYnA4TU37",
        neshan: "https://nshn.ir/3f_bvvw15xu2SY",
        embed:
          "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3236.213748866047!2d51.436706777150135!3d35.7946855237313!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x3f8e05e6a35dec77%3A0xd1e7f7e3394e6cd0!2z2LPYp9mE2YYg2LnZgtivINmG2YLYsdmH!5e0!3m2!1sen!2sat!4v1714940350017!5m2!1sen!2sat",
      },
    },
  },
};

export const title = `${wedding.couple.wife.name} & ${wedding.couple.husband.name} Wedding`;
