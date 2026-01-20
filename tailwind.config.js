/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./web/templates/**/*.templ",
    "./internal/handlers/**/*.go",
  ],
  theme: {
    extend: {
      colors: {
        // 90's Tron Grid - Blue + Purple + Cyan + Green
        neon: {
          blue: '#00D4FF',      // Electric blue - borders, grid lines
          purple: '#B026FF',    // Bright purple - headers, titles
          cyan: '#00FFFF',      // Cyan - highlights, accents
          green: '#39FF14',     // Lime green - vote scores (glowing!)
          darkPurple: '#6A0DAD', // Deep purple - secondary
          teal: '#00CED1',      // Teal - alternative
        },
      },
      fontFamily: {
        pixel: ['"Press Start 2P"', 'cursive'],
        retro: ['"VT323"', 'monospace'],
      },
      animation: {
        'glow': 'glow 2s ease-in-out infinite alternate',
        'blink': 'blink 1s linear infinite',
      },
      keyframes: {
        glow: {
          '0%': {
            textShadow: '0 0 5px #FF10F0, 0 0 10px #FF10F0, 0 0 15px #FF10F0',
          },
          '100%': {
            textShadow: '0 0 10px #FF10F0, 0 0 20px #FF10F0, 0 0 30px #FF10F0',
          },
        },
        blink: {
          '0%, 49%': { opacity: '1' },
          '50%, 100%': { opacity: '0' },
        },
      },
    },
  },
  plugins: [],
}
