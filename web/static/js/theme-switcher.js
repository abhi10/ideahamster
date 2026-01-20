// Theme Switcher - localStorage + DOM updates
(function() {
  const THEME_KEY = 'idea-hamster-theme';
  const DEFAULT_THEME = 'arctic';

  const themes = {
    'original': 'Original Pink',
    'arctic': 'Arctic Vapor',
    'rivian-earth': 'Rivian Earth',
    'rivian-cyber': 'Rivian Cyber'
  };

  // Get saved theme or default
  function getSavedTheme() {
    return localStorage.getItem(THEME_KEY) || DEFAULT_THEME;
  }

  // Save theme to localStorage
  function saveTheme(themeName) {
    localStorage.setItem(THEME_KEY, themeName);
  }

  // Apply theme to document
  function applyTheme(themeName) {
    document.documentElement.setAttribute('data-theme', themeName);

    // Update theme selector if exists
    const selector = document.getElementById('theme-selector');
    if (selector) {
      selector.value = themeName;
    }

    // Update theme button text if exists
    const themeButton = document.getElementById('current-theme-name');
    if (themeButton && themes[themeName]) {
      themeButton.textContent = themes[themeName];
    }
  }

  // Switch to a new theme
  function switchTheme(themeName) {
    if (themes[themeName]) {
      applyTheme(themeName);
      saveTheme(themeName);
    }
  }

  // Initialize theme on page load
  function initTheme() {
    const savedTheme = getSavedTheme();
    applyTheme(savedTheme);
  }

  // Expose functions globally
  window.ThemeSwitcher = {
    switch: switchTheme,
    init: initTheme,
    current: getSavedTheme,
    themes: themes
  };

  // Auto-initialize on DOM ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initTheme);
  } else {
    initTheme();
  }
})();
