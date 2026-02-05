/**
 * Edit Business 插件 API 配置模板
 * 
 * 使用说明：
 * 1. 复制此文件为 api-config.js
 * 2. 根据部署环境选择 BASE_URL
 * 3. 确保域名与实际部署地址一致
 */

const API_CONFIG = {
    // ============================================
    // 环境选择（取消注释要使用的环境）
    // ============================================
    
    // 开发环境（本地开发）
    BASE_URL: 'http://localhost:8084',

    // 生产环境（HTTPS）
    // BASE_URL: 'https://edit.crazyaigc.com',

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
    },

    /**
     * 同步单篇笔记
     * @param {Object} noteData - 笔记数据
     * @returns {Promise} 同步结果
     */
    async syncNote(noteData) {
        const url = this.getUrl('SYNC_SINGLE_NOTE');
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(noteData),
        });
        return response.json();
    },

    /**
     * 批量同步笔记
     * @param {Array} notesData - 笔记数组
     * @returns {Promise} 同步结果
     */
    async syncBatchNotes(notesData) {
        const url = this.getUrl('SYNC_BLOGGER_NOTES');
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(notesData),
        });
        return response.json();
    },

    /**
     * 同步博主信息
     * @param {Object} bloggerData - 博主数据
     * @returns {Promise} 同步结果
     */
    async syncBlogger(bloggerData) {
        const url = this.getUrl('SYNC_BLOGGER_INFO');
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(bloggerData),
        });
        return response.json();
    },
};
