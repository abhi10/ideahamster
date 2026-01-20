/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./web/templates/**/*.templ",
    "./internal/handlers/**/*.go",
  ],
  theme: {
    extend: {
      colors: {
        // Arctic Vapor - Cyan + Purple + Ice Blue + White
        neon: {
          cyan: '#00FFFF',       // Bright cyan - primary, headers
          purple: '#9D00FF',     // Deep purple - accents, titles
          ice: '#E0F7FF',        // Ice blue - highlights, soft glow
          white: '#FFFFFF',      // Pure white - text highlights
          lavender: '#E6E6FA',   // Lavender - subtle accents
          violet: '#8A2BE2',     // Blue violet - secondary purple
          aqua: '#00CED1',       // Aqua - alternative cyan
          mint: '#98FF98',       // Mint green - success states
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
