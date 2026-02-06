import { defineConfig } from 'vitest/config'
import { svelte } from '@sveltejs/vite-plugin-svelte'

export default defineConfig({
	plugins: [svelte({ hot: !process.env.VITEST })],
	test: {
		globals: true,
		environment: 'jsdom',
		setupFiles: ['./src/test/setup.ts'],
		coverage: {
			provider: 'v8',
			reporter: ['text', 'json', 'html'],
			exclude: [
				'node_modules/',
				'src/test/setup.ts',
				'*.test.ts',
				'*.test.svelte',
				'src-tauri/',
				'build/',
				'dist/',
			],
			lines: 80,
			functions: 80,
			branches: 80,
			statements: 80,
		},
		include: ['src/**/*.{test,spec}.{ts,js,svelte}'],
		exclude: ['node_modules/', 'dist/', 'build/', 'src-tauri/'],
	},
})
