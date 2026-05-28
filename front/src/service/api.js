/**
 * - 接口地址统一管理
 * - 以对象形式导出
 * - 属性包含当前项目“所有用到的接口”
 */
const BASE_URL = import.meta.env.PROD ? location.pathname.replace(/\/+$/, '') : "/api/nginx";

export default {
    host: `${BASE_URL}`,
    health: `${BASE_URL}/health`,
    metrics: `${BASE_URL}/metrics`
}

