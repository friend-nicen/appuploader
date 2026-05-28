/**
 * @author 友人a丶
 * @date 2026-04-11
 * 弹窗打开时禁止页面滚动，关闭时恢复
 * 支持嵌套弹窗计数，不会互相覆盖
 */
import {onUnmounted, watch} from 'vue';

/* 全局计数器，解决多层弹窗嵌套问题 */
let lockCount = 0;
const lockClass = 'with_overflow';

export default function useBodyScrollLock(props, type = 1) {

    /* 监听的属性名 */
    const propName = type === 1 ? 'value' : 'visible';

    /* 锁定滚动 */
    function lockBody() {
        lockCount++;
        if (lockCount === 1) {
            document.body.classList.add(lockClass);
        }
    }

    /* 解锁滚动 */
    function unlockBody() {
        lockCount = Math.max(0, lockCount - 1);
        if (lockCount === 0) {
            document.body.classList.remove(lockClass);
        }
    }

    /* 强制清空锁定 */
    function forceUnlockBody() {
        lockCount = 0;
        document.body.classList.remove(lockClass);
    }

    /* 初始化时执行一次 */
    const unwatch = watch(
        () => props[propName],
        (visible) => {
            if (visible) {
                lockBody();
            } else {
                unlockBody();
            }
        },
        {immediate: true}
    );

    /* 组件销毁时强制解锁，避免页面卡死 */
    onUnmounted(() => {
        unwatch();
        forceUnlockBody();
    });
}