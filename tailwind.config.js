/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    colors: {
      transparent: 'transparent',
      current: 'currentColor',
      'white': '#ffffff',
      'arisu': {
        100: '#ecf1f5',
        200: '#c7d4e2',
        300: '#a2b8cf',
        400: '#7c9bbb',
        500: '#577ea8',
        600: '#446283',
        700: '#30465d',
        800: '#1d2a38',
        900: '#0a0e13',
      },
      // ...
    },
  },
  plugins: [],
}