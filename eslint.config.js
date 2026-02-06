import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';
import svelteESLint from 'eslint-plugin-svelte';
import globals from 'globals';

export default [
  eslint.configs.recommended,
  ...tseslint.configs.recommended,
  ...svelteESLint.configs['flat/recommended'],
  {
    languageOptions: {
      ecmaVersion: 2020,
      sourceType: 'module',
      globals: {
        ...globals.browser,
        ...globals.node,
      },
    },
  },
  {
    ignores: [
      'build/',
      '.svelte-kit/',
      'dist/',
      'node_modules/',
      'src-tauri/target/',
    ],
  },
  prettier.configs.recommended,
];
