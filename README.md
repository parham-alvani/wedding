# Wedding RSVP Website

<div align="center">

![Banner](docs/banner.webp)

</div>

## Introduction

We (Elaheh and Parham) are going to get married.
This repository will provide the invitation card and details of the ceremony.

> [!note]
> This project is also designed as a **template** for anyone who wants to create their own wedding RSVP website. Fork it and customize it for your own wedding!

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
- Boho Floral Design Studio - 2
- Kourosh Ceremonial Services
- Aida Mokhtary
- Ajoodaniyeh Mansion
- Makan Mansion
- Parastoo Mezon
- Saadat Rent - 1
- Saadat Rent - 2

## Running the Project

This project uses [just](https://github.com/casey/just) as a command runner. Install it first, then:

```bash
# Install all dependencies
just install

# Start both frontend and backend for development
just dev

# Build everything for production
just build

# Run all tests
just test
```

### Backend (WedBack)

The backend is a Go application that manages guests and serves the API on `:1378`.

```bash
just back serve       # Build and start the server
just back insert      # Add a new guest (interactive TUI)
just back list        # List all guests with RSVP status
just back test        # Run tests
just back lint        # Run linter
```

Configuration can be provided via `config.toml` or environment variables prefixed with `wedback_`.

#### Adding a new guest

`just back insert` opens an interactive TUI where you fill in:

1. First Name
2. Last Name
3. Partner's First Name
4. Partner's Last Name
5. Is Family? (`true` / `false`)
6. Children (number)

Use `Tab` / `Shift+Tab` to navigate between fields and `Enter` to confirm.

### Frontend (WedFront)

<div align="center">

[![Built with Astro](https://astro.badg.es/v2/built-with-astro/small.svg)](https://astro.build)

</div>

The frontend runs as a standalone Astro Node.js server (SSR). It proxies API requests to the backend internally.

```bash
just front install    # Install pnpm dependencies
just front dev        # Start dev server at localhost:4321
just front build      # Build for production
just front serve      # Start production server
just front format     # Format code with Prettier
just front clean      # Remove node_modules and build output
```

Set `WEDFRONT_BACKEND_URL` to point to the backend (defaults to `http://127.0.0.1:1378`).

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
