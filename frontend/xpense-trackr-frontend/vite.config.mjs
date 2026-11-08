import { defineConfig, transformWithEsbuild } from 'vite';
import react from '@vitejs/plugin-react';

const jsxInJsPlugin = {
    name: 'jsx-in-js',
    enforce: 'pre',
    async transform(code, id) {
        if (!/\/src\/.*\.js$/.test(id)) return null;
        return transformWithEsbuild(code, id, { loader: 'jsx', jsx: 'automatic' });
    },
};

export default defineConfig({
    plugins: [jsxInJsPlugin, react()],
    define: {
        'process.env.REACT_APP_API_URL': JSON.stringify(
            process.env.REACT_APP_API_URL || process.env.VITE_API_URL || '',
        ),
    },
    server: {
        host: '0.0.0.0',
        port: 3000,
        strictPort: true,
    },
    build: {
        outDir: 'build',
        emptyOutDir: true,
    },
});
