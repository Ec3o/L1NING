module.exports = {
  devServer: {
    proxy: {
      '/api': {
        target: 'http://localhost:8090',
        changeOrigin: true,
        pathRewrite: { '^/api': '/api' },
      },
    },
    historyApiFallback: true,
  },
  transpileDependencies: [],
};
