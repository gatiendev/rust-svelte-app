```markdown
# SvelteKit Dark Starter Template

A production-ready SvelteKit starter template featuring a dark theme, modular UI components, authentication UI with mock state, and a robust testing setup. Built with TypeScript, TailwindCSS (v3 via PostCSS), and designed for scalability and developer experience.

## 🚀 Quick Start

```bash
# Clone the repository (or use degit)
npx degit your-username/sveltekit-dark-starter my-app
cd my-app

# Install dependencies
npm install

# Start development server
npm run dev

# Visit http://localhost:5173
```

## 📁 Project Architecture

```
.
├── .gitignore
├── package.json
├── postcss.config.js          # Tailwind + Autoprefixer config
├── tailwind.config.js          # Tailwind theme & dark mode config
├── svelte.config.js            # SvelteKit adapter & preprocessors
├── vite.config.ts              # Vite + Vitest configuration
├── README.md
└── src/
    ├── app.html                # HTML entry with inline dark mode script
    ├── app.css                 # Tailwind imports & custom theme variables
    ├── lib/
    │   ├── components/
    │   │   ├── layout/
    │   │   │   ├── Header.svelte    # Responsive header with theme toggle
    │   │   │   └── Footer.svelte    # Simple footer
    │   │   ├── ui/                  # Reusable UI components
    │   │   │   ├── Button.svelte
    │   │   │   ├── Card.svelte
    │   │   │   ├── Input.svelte
    │   │   │   ├── Modal.svelte
    │   │   │   └── Alert.svelte
    │   │   ├── auth/                 # Authentication forms
    │   │   │   ├── LoginForm.svelte
    │   │   │   └── RegisterForm.svelte
    │   │   └── shared/                # Placeholder for shared utilities
    │   ├── stores/                    # Global state
    │   │   ├── theme.ts                # Dark mode store with localStorage
    │   │   └── auth.ts                 # Mock authentication store
    │   └── types/                       # TypeScript type definitions
    └── routes/
        ├── +layout.svelte               # Global layout (Header, Footer, main)
        ├── +page.svelte                  # Landing page with hero & features
        ├── auth/
        │   ├── login/
        │   │   └── +page.svelte          # Login page
        │   └── register/
        │       └── +page.svelte          # Register page
        └── __tests__/                     # Unit tests
            └── Button.test.ts
```

## 🎨 Design System & Dark Mode

- **Dark mode** is implemented using Tailwind's `class` strategy. The theme store (`stores/theme.ts`) toggles the `dark` class on the `<html>` element and persists the preference in `localStorage`.
- To prevent a flash of light mode on initial load, an inline script in `app.html` reads `localStorage` and applies the correct class before the page paints.
- The color palette is defined in `tailwind.config.js` under the `primary` key. You can easily change the accent colors by modifying this config.

### Theme Store (`stores/theme.ts`)

```typescript
import { writable } from 'svelte/store';
import { browser } from '$app/environment';

type Theme = 'light' | 'dark';

function createThemeStore() {
  const stored = browser ? (localStorage.getItem('theme') as Theme) : 'dark';
  const initial = stored === 'light' ? 'light' : 'dark';
  const { subscribe, update } = writable(initial);

  if (browser) {
    subscribe((value) => {
      localStorage.setItem('theme', value);
      document.documentElement.classList.toggle('dark', value === 'dark');
    });
  }

  return {
    subscribe,
    toggle: () => update(t => t === 'light' ? 'dark' : 'light'),
  };
}

export const theme = createThemeStore();
export const toggleTheme = () => theme.toggle();
```

## 🔐 Authentication (Mock)

Authentication is handled via a Svelte store (`stores/auth.ts`) with mock implementations of `login`, `register`, and `logout`. It simulates API delays and basic validation. The store exposes:

- `user`: `User | null`
- `loading`: `boolean`
- `error`: `string | null`
- `login(email, password)`
- `register(email, password, name?)`
- `logout()`

The forms (`LoginForm.svelte`, `RegisterForm.svelte`) use the reusable UI components and dispatch `submit` events. The parent pages listen to these events and call the appropriate auth methods.

## 🧩 UI Components

All components are located in `src/lib/components/ui/` and are designed to be highly reusable, testable, and accessible. They accept props, dispatch events, and support dark mode via Tailwind's `dark:` variants.

### Button (`Button.svelte`)

- **Variants**: `primary`, `secondary`, `outline`, `ghost`
- **Sizes**: `sm`, `md`, `lg`
- **Props**: `disabled`, `fullWidth`, `type`
- **Events**: forwards native `click` event (no duplicate handling)
- **Styles**: Smooth transitions, hover scaling, focus rings

### Card (`Card.svelte`)

- **Props**: `padding`, `shadow` ('none'|'sm'|'md'|'lg'), `border`, `hover`, `interactive`
- **Slots**: default content
- **Styles**: Rounded corners, background, transitions

### Input (`Input.svelte`)

- **Props**: `label`, `type`, `value`, `placeholder`, `error`, `disabled`, `id`, `icon` (component)
- **Events**: binds to `value`
- **Styles**: Focus states, error styling, icon positioning

### Modal (`Modal.svelte`)

- **Props**: `open`, `title`, `closeOnOutsideClick`, `size` ('sm'|'md'|'lg')
- **Events**: `close` when modal is closed
- **Transitions**: `fade` on backdrop, `fly` on modal content
- **Styles**: Backdrop blur, shadow, rounded corners

### Alert (`Alert.svelte`)

- **Props**: `type` ('success'|'error'|'info'|'warning'), `dismissible`, `title`
- **Events**: `dismiss` when close button clicked
- **Styles**: Color-coded backgrounds, icons, borders

## 🧪 Testing

The project is configured with **Vitest** and **Testing Library** for unit testing Svelte components.

- Test files are co-located in `__tests__` folders (e.g., `routes/__tests__/Button.test.ts`).
- Run tests: `npm run test`
- Run with UI: `npm run test:ui`

Example test (`Button.test.ts`):

```typescript
import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import Button from '$lib/components/ui/Button.svelte';

describe('Button', () => {
  it('renders with default props', () => {
    render(Button, { props: { children: 'Click me' } });
    expect(screen.getByRole('button')).toHaveClass('bg-primary-600');
  });

  it('handles click events', async () => {
    const onClick = vi.fn();
    render(Button, { props: { onclick: onClick } });
    await userEvent.click(screen.getByRole('button'));
    expect(onClick).toHaveBeenCalledOnce();
  });
});
```

## 🛠 Customization Guide

### Changing the Accent Color

Edit the `primary` color palette in `tailwind.config.js`:

```js
theme: {
  extend: {
    colors: {
      primary: {
        50: '#your-color',
        600: '#your-color',
        // ... all shades
      },
    },
  },
},
```

Then replace any hardcoded color classes (e.g., `bg-blue-600`) with `bg-primary-600` in your components.

### Adding New UI Components

1. Create a new `.svelte` file in `src/lib/components/ui/`.
2. Follow the existing patterns: use TypeScript, define props with `export let`, and style with Tailwind.
3. Add dark mode variants using `dark:` prefix.
4. Write a test in a `__tests__` folder.

### Modifying the Layout

- Edit `Header.svelte` and `Footer.svelte` in `src/lib/components/layout/`.
- Update `+layout.svelte` to change the overall page structure.

### Authentication Backend Integration

Replace the mock auth store (`stores/auth.ts`) with real API calls. The `login` and `register` methods should call your backend endpoints. Keep the same interface to ensure forms work without changes.

## 🌐 Deployment

This template is configured with `@sveltejs/adapter-auto`, which automatically selects the appropriate adapter for your deployment platform (Vercel, Netlify, Cloudflare, etc.).

To build for production:

```bash
npm run build
```

The output will be in the `build` directory (or platform-specific folder).

## 📚 Learn More

- [SvelteKit Documentation](https://kit.svelte.dev/docs)
- [TailwindCSS Documentation](https://tailwindcss.com/docs)
- [Testing Library for Svelte](https://testing-library.com/docs/svelte-testing-library/intro)

## 🤝 Contributing

Feel free to submit issues or pull requests to improve this template. For major changes, please open an issue first to discuss what you would like to change.

## 📄 License

MIT

```

This README provides a comprehensive overview for another LLM to understand the project architecture, key files, and how to extend it. It includes setup instructions, component documentation, testing details, and customization guidance.
