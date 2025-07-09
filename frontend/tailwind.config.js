/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'splunk-dark': '#1a1a1a',
        'splunk-darker': '#0f0f0f',
        'splunk-gray': '#2d2d2d',
        'splunk-light-gray': '#404040',
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
} 