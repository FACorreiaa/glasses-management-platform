// tailwind.config.mjs

import forms from '@tailwindcss/forms';
import typography from '@tailwindcss/typography';
import daisyui from 'daisyui';

/** @type {import('tailwindcss').Config} */
export default { // Changed from module.exports
  content: ['./app/**/*.{css,view,js,templ,html}'], // Ensure this covers all files using Tailwind classes
  theme: {
    container: {
      center: true,
      padding: "2rem",
      screens: {
        "2xl": "1400px"
      }
    },
    // NOTE: Defining `colors` outside `extend` REPLACES Tailwind's default color palette.
    // Keep this only if you explicitly want to remove all default Tailwind colors.
    // It's usually better to put ALL color definitions inside `extend`.
    colors: {
      blue: '#1fb6ff',
      purple: '#7e5bef',
      pink: '#ff49db',
      orange: '#ff7849',
      green: '#13ce66',
      yellow: '#ffc82c',
      'gray-dark': '#273444',
      gray: '#8492a6',
      'gray-light': '#d3dce6',
    },
    fontFamily: {
      sans: ['Graphik', 'sans-serif'],
      serif: ['Merriweather', 'serif'],
      lato: ['Lato', 'sans-serif'],
    },
    extend: { // Use extend to ADD to Tailwind's defaults
      colors: { // These colors are ADDED to Tailwind's defaults (or your overrides above)
        border: "hsl(var(--border) / <alpha-value>)",
        input: "hsl(var(--input) / <alpha-value>)",
        ring: "hsl(var(--ring) / <alpha-value>)",
        background: "hsl(var(--background) / <alpha-value>)",
        foreground: "hsl(var(--foreground) / <alpha-value>)",
        primary: {
          DEFAULT: "hsl(var(--primary) / <alpha-value>)",
          foreground: "hsl(var(--primary-foreground) / <alpha-value>)"
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary) / <alpha-value>)",
          foreground: "hsl(var(--secondary-foreground) / <alpha-value>)"
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive) / <alpha-value>)",
          foreground: "hsl(var(--destructive-foreground) / <alpha-value>)"
        },
        muted: {
          DEFAULT: "hsl(var(--muted) / <alpha-value>)",
          foreground: "hsl(var(--muted-foreground) / <alpha-value>)"
        },
        accent: {
          DEFAULT: "hsl(var(--accent) / <alpha-value>)",
          foreground: "hsl(var(--accent-foreground) / <alpha-value>)"
        },
        popover: {
          DEFAULT: "hsl(var(--popover) / <alpha-value>)",
          foreground: "hsl(var(--popover-foreground) / <alpha-value>)"
        },
        card: {
          DEFAULT: "hsl(var(--card) / <alpha-value>)",
          foreground: "hsl(var(--card-foreground) / <alpha-value>)"
        }
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
        '4xl': '2rem',

      },
      spacing: {
        '8xl': '96rem',
        '9xl': '128rem',
      },
    },
  },
  daisyui: {
    themes: [
      'light',
      'dark',
      'nord',
      'cyberpunk',
      'pastel',
      'cupcake',
      'night',
      'bumblebee',
      'business',
      'lemonade',
      // Your custom Catppuccin themes are correctly defined here
      {
        'catppuccin-latte': { /* ... */ },
        'catppuccin-frappe': { /* ... */ },
        'catppuccin-macchiato': { /* ... */ },
        'catppuccin-mocha': { /* ... */ },
      },
    ],
  },
  plugins: [
    forms,        // Use the imported variable
    typography,   // Use the imported variable
    daisyui,      // Use the imported variable
    // require('autoprefixer'), // REMOVED - Belongs in postcss.config.cjs
  ],
};