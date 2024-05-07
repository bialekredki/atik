/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/templates/*.templ", "./web/templates/components/*.templ"],

  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
}

