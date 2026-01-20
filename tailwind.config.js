/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./web/templates/**/*.templ",
    "./internal/handlers/**/*.go",
  ],
  theme: {
    extend: {
      colors: {
        // 90's retro color palette - Electric Blue + Coral + Purple
        neon: {
          blue: '#00D4FF',      // Electric blue - primary
          coral: '#FF6B6B',     // Coral - accent
          purple: '#9D00FF',    // Purple - secondary
          teal: '#00CED1',      // Teal - alternative
          gold: '#FFD700',      // Gold - success states
          green: '#39FF14',     // Lime green - kept for eligible badges
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
