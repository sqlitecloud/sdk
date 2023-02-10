const defaultTheme = require('tailwindcss/defaultTheme')
module.exports = {
  content: [
    "./app/template/**/*.{html,js}",
    "./app/src/**/*.{html,js}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter var', ...defaultTheme.fontFamily.sans]
      },
    }
  },
  plugins: [
    require('@tailwindcss/forms')
  ],
}
