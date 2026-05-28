import {defineConfig} from 'vite'
import path from 'node:path'
import vue from '@vitejs/plugin-vue'
import postcssPxtorem from 'postcss-pxtorem'
import Components from 'unplugin-vue-components/vite'
import {AntDesignVueResolver} from 'unplugin-vue-components/resolvers'

// https://vite.dev/config/
export default defineConfig(({mode}) => {
    const isDev = mode === 'development'
    return {
        plugins: [
            vue(),
            Components({
                resolvers: [
                    AntDesignVueResolver({
                        importStyle: false
                    })
                ],
            }),
        ],

        define: {
            __BUILD_TIME__: JSON.stringify(new Date().toISOString()),
        },

        resolve: {
            alias: {
                '@': path.resolve(__dirname, 'src'),
                '#': path.resolve(__dirname, 'wailsjs'),
            },
        },

        esbuild: {
            drop: !isDev ? ['console', 'debugger'] : [],
        },

        css: {
            preprocessorOptions: {
                less: {
                    javascriptEnabled: true,
                },
                scss: {
                    additionalData: `
                    @use "@/theme/scroll.scss" as *;
                    @use "@/theme/theme.scss" as *;
                    `,
                },
            },
            postcss: {
                plugins: [
                    postcssPxtorem({
                        rootValue: 16,
                        propList: [
                            'font-size',
                            'padding*',
                            'margin*',
                            '*width',
                            'height',
                            'top',
                            'left',
                            'right',
                            'bottom',
                            'line-height',
                        ],
                        selectorBlackList: ['media'],
                        unitPrecision: 8,
                    }),
                ],
            },
        },

        server: {
            host: '0.0.0.0',
            port: 8080
        },

        build: {
            sourcemap: false,
            outDir: 'dist',
            cssTarget: 'chrome50',
            rollupOptions: {
                output: {
                    chunkFileNames: 'assets/js/[hash].js',
                    entryFileNames: 'assets/js/[hash].js',
                    assetFileNames: 'assets/[ext]/[hash].[ext]',
                },
            },
        },
    }
})
