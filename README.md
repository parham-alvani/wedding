# Wedding RSVP Website

<div align="center">

![Banner](docs/banner.webp)

</div>

## Introduction

We (Elaheh and Parham) are going to get married.
This repository will provide the invitation card and details of the ceremony.

This project is also designed as a **template** for anyone who wants to create their own wedding RSVP website. Fork it and customize it for your own wedding!

## How does it happen?

The following are the tasks and things we are doing for our wedding.

- Mina Rad - 1
- Mina Rad - 2
- Hasti Mezon
- Shemiranat health and treatment network
- Farzad Ahmadi
- Rozzet Studio
- Noghreh Wedding
- Rajian
- Boho Floral Design Studio - 1
- Kourosh Ceremonial Services
- Aida Mokhtary
- Ajoodaniyeh Mansion
- Makan Mansion
- Parastoo Mezon
- Saadat Rent - 1
- Saadat Rent - 2

## Running the Project

### Backend (WedBack)

The backend is a Go application that manages guests and serves the API.

```bash
cd wedback
go build -o wedback ./cmd/wedback
```

#### Start the server

```bash
./wedback serve
```

The server listens on `:1378` by default. Configuration can be provided via `config.toml` or environment variables prefixed with `wedback_`.

#### Add a new guest

```bash
./wedback insert
```

This opens an interactive TUI where you fill in:

1. First Name
2. Last Name
3. Partner's First Name
4. Partner's Last Name
5. Is Family? (`true` / `false`)
6. Children (number)

Use `Tab` / `Shift+Tab` to navigate between fields and `Enter` to confirm.

#### List all guests

```bash
./wedback list
```

Displays a table of all guests with their RSVP status.

### Frontend (WedFront)

<div align="center">

[![Built with Astro](https://astro.badg.es/v2/built-with-astro/small.svg)](https://astro.build)

</div>

```bash
cd wedfront
npm install
```

| Command                | Action                                            |
| :--------------------- | :------------------------------------------------ |
| `npm run dev`          | Start local dev server at `localhost:4321`        |
| `npm run build`        | Build your production site to `./dist/`           |
| `npm run preview`      | Preview your build locally, before deploying      |
| `npm run astro ...`    | Run CLI commands like `astro add`, `astro check`  |
| `npm run astro --help` | Get help using the Astro CLI                      |
| `npm run format`       | Format code with [Prettier](https://prettier.io/) |
| `npm run clean`        | Remove `node_modules` and build output            |

The frontend expects the backend to be running. Set `WEDFRONT_BACKEND_URL` to point to the backend (defaults to `http://127.0.0.1:1378`).

## Customization

To use this project for your own wedding, fork the repository and edit the following files:

### Frontend

Edit `wedfront/src/wedding.config.ts` to set your names, dates, social links, music, and site URLs:

```typescript
export const wedding = {
  couple: {
    husband: { name: "...", namePersian: "...", lastName: "...", emoji: "...", socials: { ... } },
    wife:    { name: "...", namePersian: "...", lastName: "...", emoji: "...", socials: { ... } },
    lastNamePersian: "...",
  },
  dates: {
    wedding: "Jun 16, 2024 18:30:00+03:30",
    engaged: "May 10, 2024 19:00:00+03:00",
  },
  site: { url: "...", github: "..." },
  music: { ceremony: "...", engaged: "..." },
};
```

All pages and components read from this single config file.

### Backend

Edit `wedback/config.toml` or use environment variables to set:

```toml
[wedding]
husband_name = "Your Name"
wife_name = "Partner Name"
base_url = "https://your-wedding-site.com"
```

Or via environment variables: `wedback_wedding__husband_name`, `wedback_wedding__wife_name`, `wedback_wedding__base_url`.
