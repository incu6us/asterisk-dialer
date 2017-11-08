const webpack = require('webpack');
const path = require('path');
const HtmlWebpackPlugin  = require('html-webpack-plugin');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const ProgressBarPlugin = require('progress-bar-webpack-plugin');
const WebpackNotifierPlugin = require('webpack-notifier');


module.exports = {
    entry: ['babel-polyfill',path.resolve(__dirname + '/src/index.js')],
    output: {
        path: path.resolve(__dirname + '/build'),
        filename: 'static/js/[name].[hash:8].js',
    },
    module: {
        rules: [
            {
                test: /\.js$/,
                loader: 'babel-loader',
                query: {
                    babelrc: false,
                    presets: ['babel-preset-react', 'env', 'babel-preset-stage-2'],
                    minified: true
                },
                exclude: /node_modules/
            },
            {
                test: /\.scss$/,
                use: ExtractTextPlugin.extract({
                    fallback: 'style-loader',
                    use: [
                        {
                            loader: 'css-loader',
                        },
                        'sass-loader'
                    ]
                })
            }
        ]
    },
    plugins: [
        // Generates an `index.html` file with the <script> injected.
        new HtmlWebpackPlugin({
            inject: true,
            template: path.resolve(__dirname + '/index.html')
        }),
        new webpack.DefinePlugin({
            'process.env': {
                NODE_ENV: JSON.stringify(process.env.NODE_ENV)
            }
        }),
        new WebpackNotifierPlugin({
            alwaysNotify: true
        }),
        new webpack.HotModuleReplacementPlugin(),
        new ProgressBarPlugin(),
        new ExtractTextPlugin('app.css')
    ].concat(
        process.env.NODE_ENV === 'production'
            ? [
                new webpack.optimize.UglifyJsPlugin({
                    compress: {warnings: false, comparisons: false},
                    output: {comments: false, ascii_only: true}
                })
            ] : []
    ),
    resolve: {
        modules: [
            path.resolve(__dirname, 'node_modules'),
        ],
        extensions: ['.js', '.jsx', '.scss', '.json'],
    },
    devServer: {
        port: 3030
    },
    devtool: process.env.NODE_ENV === 'production' ? 'source-map': 'eval'
};
