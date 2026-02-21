import * as js from '@eslint/js'
import * as globals from 'globals'
import tseslint from 'typescript-eslint'
import * as pluginVue from 'eslint-plugin-vue'
import { defineConfig } from 'eslint/config'
import * as vueParser from 'vue-eslint-parser'

export default defineConfig([
  {
    ignores: [
      'dist/**',
      'node_modules/**',
    ],
  },
  {
    files: ['**/*.{ts,vue}'],
    plugins: { js },
    extends: ['js/recommended'],
    languageOptions: { globals: globals.browser },
  },
  tseslint.configs.recommended,
  // tseslint.configs.strictTypeChecked,
  pluginVue.configs['flat/recommended-error'],
  // pluginVue.configs["flat/recommended-error"],
  {
    files: ['**/*.vue'],
    languageOptions: {
      parser: vueParser,
      parserOptions: {
        parser: tseslint.parser,
        extraFileExtensions: ['.vue'],
      },
    },
    rules: {
      'vue/component-name-in-template-casing': ['error', 'PascalCase'],
      'vue/multi-word-component-names': 'off',
    },
  },
  {
    rules: {
      '@typescript-eslint/consistent-type-imports': ['error', { prefer: 'type-imports' }],
      semi: ['error', 'never'],
      indent: ['error', 2],
      quotes: ['error', 'single'],
      'comma-dangle': ['error', 'always-multiline'],
      'object-curly-spacing': ['error', 'always'],
      'eqeqeq': ['error', 'always'],
      'no-debugger': 'error',
      'no-console': ['error', { allow: ['warn', 'error'] }],
      'vue/no-mutating-props': 'error',
    },
  },
])

