// XHS2Feishu 后端API配置
// 请在这里配置你的后端服务器地址

const API_CONFIG = {
    // 后端服务器地址
    // 生产环境：使用Dokploy服务器地址
    BASE_URL: 'https://api.keenchase.com',

    // 开发环境：使用本地地址（如果需要本地开发，取消下面的注释）
    // BASE_URL: 'http://localhost:3000',

    // API端点
    ENDPOINTS: {
        SYNC_SINGLE_NOTE: '/api/sync/single-note',
        SYNC_BLOGGER_NOTES: '/api/sync/blogger-notes',
        SYNC_BLOGGER_INFO: '/api/sync/blogger-info',
        VERIFY_ORDER: '/api/verify/order'
    }
};

// 导出配置
if (typeof module !== 'undefined' && module.exports) {
    module.exports = API_CONFIG;
}
