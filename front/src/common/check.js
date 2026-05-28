import user from '@/stores/user'

/**
 * 用户权限校验
 * @param item
 * @return {boolean}
 */
export function check(item) {

    /* 不需要权限控制 */
    if (!item.permiss) {
        return true;
    }

    /** @type Suser 获取用户信息 */
    const userInfo = user();

    /* 针对数组和字符串进行不同的判断 */
    if (Array.isArray(item.permiss)) {
        for (let i of item.permiss) {
            if (userInfo.permissions.indexOf(i) === -1) {
                return false;
            }
        }
        return true;
    } else {
        return userInfo.permissions.indexOf(item.permiss) > -1;
    }

}