/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./src/**/*.{html,js,svelte,ts}",
    "./src/lib/**/*.{html,js,svelte,ts}"
  ],
  theme: {
    extend: {
      colors: {
        'divoom-primary': '#FF6B35',
        'divoom-secondary': '#004E89',
        'divoom-dark': '#1A1A2E',
        'divoom-light': '#F7F7F7'
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
