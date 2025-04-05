module.exports = {
  plugins: {
    'postcss-import': {},
    // Note: Consider if postcss-url should run before or after tailwind
    'postcss-url': {
      url: 'copy',
      useHash: true,
      assetsPath: '../fonts', // Make sure this path is correct relative to output.css
    },
    '@tailwindcss/postcss': {}, // <-- Use the new package name as a key
    'autoprefixer': {},
  }
};