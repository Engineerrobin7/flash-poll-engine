# 🎨 Flash-Poll Frontend (React)

A high-contrast, Neo-Brutalist interface built for speed and clarity.

## Features

- **Reactive State**: Custom `usePolls` hook handles both initial fetch and live SSE synchronization.
- **Zero Local Persistence**: Strictly follows the requirement to not store votes in LocalStorage.
- **Tactile UI**: High-vis borders, "sinking" button states, and CSS transitions that emphasize movement.
- **Live Sync**: Uses a native `EventSource` connection to the Go backend for real-time bar updates.

## Setup

```bash
npm install
npm run dev
```

## Styling Notes
The UI uses **CSS Variables** for easy theme switching. Categories like `Tech`, `Business`, and `Design` have unique color tokens defined in `index.css`.

---
*Built with React 19 + Vite.*
