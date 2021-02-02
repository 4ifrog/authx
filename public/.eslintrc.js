module.exports = {
  extends: ['plugin:prettier/recommended', 'prettier/babel'],
  env: {
    browser: true,
    es6: true,
  },
  parserOptions: {
    ecmaVersion: 2018,
    sourceType: 'module',
  },
  rules: {
    // Add or overwrite specific rules.
    'arrow-parens': ['error'],
    'import/prefer-default-export': ['off'],
    'no-param-reassign': ['warn'],
    'no-unused-vars': ['warn'],
  },
};
