# Webpack 构建优化配置
# 使用 React Scripts v5 的构建优化

const { overrideCraco, CracoPlugin } = require('@craco/craco');
const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');
const TerserPlugin = require('terser-webpack-plugin');
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin');
const CompressionWebpackPlugin = require('compression-webpack-plugin');

module.exports = {
    webpack: {
        optimization: {
            minimize: true,
            splitChunks: {
                chunks: 'all',
                cacheGroups: {
                    vendor: {
                        test: /[\\/]node_modules[\\/]/,
                        name: 'vendors',
                        priority: 10,
                        reuseExistingChunk: true,
                    },
                    antd: {
                        test: /[\\/]node_modules[\\/](antd|@ant-design)/,
                        name: 'antd',
                        priority: 20,
                        reuseExistingChunk: true,
                    },
                    react: {
                        test: /[\\/]node_modules[\\/](react|react-dom|react-router|react-router-dom)/,
                        name: 'react',
                        priority: 30,
                        reuseExistingChunk: true,
                    },
                    echarts: {
                        test: /[\\/]node_modules[\\/](echarts)/,
                        name: 'echarts',
                        priority: 40,
                        reuseExistingChunk: true,
                    },
                    monaco: {
                        test: /[\\/]node_modules[\\/](monaco-editor)/,
                        name: 'monaco',
                        priority: 50,
                        reuseExistingChunk: true,
                    },
                    codemirror: {
                        test: /[\\/]node_modules[\\/](codemirror)/,
                        name: 'codemirror',
                        priority: 60,
                        reuseExistingChunk: true,
                    },
                },
            },
        },
    },
    plugins: [
        new BundleAnalyzerPlugin({
            analyzerMode: 'static',
            openAnalyzer: false,
            reportFilename: '../report.html',
            defaultSizes: {
                javascript: 400,
                vendor: 300,
                common: 200,
            },
        }),
        new CompressionWebpackPlugin({
            filename: '[name].[contenthash].js.gz',
            algorithm: 'gzip',
            test: /\.(js|css|html|json)$/,
            threshold: 10240, // 大于10KB才压缩
            minRatio: 0.8,
            deleteOriginalAssets: true,
        }),
        new CssMinimizerPlugin({
            test: /\.css$/,
        }),
        new TerserPlugin({
            terserOptions: {
                ecmaVersion: 2020,
                useCache: true,
                parallel: true,
                sourceMap: false,
            },
        }),
    ],
    module: {
        rules: [
            {
                test: /\.(js|jsx|ts|tsx)$/,
                exclude: /node_modules/,
                use: {
                    loader: TerserPlugin.loader,
                    options: {
                        cacheDirectory: true,
                        cacheCompression: true,
                    },
                },
            },
        ],
    },
};