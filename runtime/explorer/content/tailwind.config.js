/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
    "./**/*.{js,ts,jsx,tsx}", // since your files are in the root
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};
