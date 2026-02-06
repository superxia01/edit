/**
 * Edit Business Chrome插件 API 配置
 *
 * 使用说明：
 * 根据部署环境选择 BASE_URL（取消注释对应行）
 */

const API_CONFIG = {
    // ============================================
    // 环境选择（取消注释要使用的环境）
    // ============================================

    // 生产环境（HTTPS）- 默认使用
    BASE_URL: 'https://edit.crazyaigc.com',

    // 开发环境（本地开发）
    // BASE_URL: 'http://localhost:8084',

    // ============================================
    // API 端点配置
    // ============================================

    ENDPOINTS: {
        // 同步单篇笔记
        SYNC_SINGLE_NOTE: '/api/v1/notes',

        // 批量同步笔记
        SYNC_BLOGGER_NOTES: '/api/v1/notes/batch',

        // 同步博主信息
        SYNC_BLOGGER_INFO: '/api/v1/bloggers',

        // 获取七牛云上传token
        QINIU_UPLOAD_TOKEN: '/api/v1/qiniu/upload-token',
    },

    // ============================================
    // 辅助方法
    // ============================================

    /**
     * 获取完整的 API URL
     * @param {string} endpoint - 端点名称
     * @returns {string} 完整的 API URL
     */
    getUrl(endpoint) {
        return `${this.BASE_URL}${this.ENDPOINTS[endpoint]}`;
    }
};

// 导出配置
if (typeof module !== 'undefined' && module.exports) {
    module.exports = API_CONFIG;
}
