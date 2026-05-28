/**
 * 自定义事件管理器
 * 使用 Map 存储事件，key 为事件名称，value 为监听器数组
 */
class EventBus {
    constructor() {
        this._event = new Map();
    }

    /**
     * 添加事件监听器
     * @param {string} eventName - 事件名称
     * @param {Function} listener - 监听器函数
     * @param {boolean} [once=false] - 是否为单次监听（触发后自动移除）
     * @returns {EventBus} 返回实例自身，支持链式调用
     */
    add(eventName, listener, once = false) {
        if (typeof listener !== 'function') {
            throw new TypeError('listener must be a function');
        }
        if (!this._event.has(eventName)) {
            this._event.set(eventName, []);
        }
        // 存储监听器及是否单次执行的标记
        this._event.get(eventName).push({fn: listener, once});
        return this;
    }

    /**
     * 添加单次事件监听器（触发后自动移除）
     * @param {string} eventName - 事件名称
     * @param {Function} listener - 监听器函数
     * @returns {EventBus} 返回实例自身，支持链式调用
     */
    once(eventName, listener) {
        return this.add(eventName, listener, true);
    }

    /**
     * 触发指定事件
     * @param {string} eventName - 事件名称
     * @param  {...any} args - 传递给监听器的参数
     * @returns {EventBus} 返回实例自身，支持链式调用
     */
    emit(eventName, ...args) {
        const listeners = this._event.get(eventName);
        if (!listeners || !listeners.length) return this;
        /* 复制一份监听器数组，避免触发过程中数组长度变化导致的问题 */
        const listenersCopy = [...listeners];
        for (const item of listenersCopy) {
            try {
                item.fn.apply(this, args);
            } catch (error) {
                console.error(`[EventBus] 监听器执行失败 - ${eventName}:`, error);
            }
            /* 如果是单次监听，执行后移除该监听器 */
            if (item.once) {
                this.remove(eventName, item.fn);
            }
        }
        return this;
    }

    /**
     * 获取指定事件名称对应值的最后一个事件监听器
     * @param {string} eventName - 事件名称
     * @returns {Function|null} 最后一个监听器，不存在则返回 null
     */
    get(eventName) {
        const listeners = this._event.get(eventName);
        return listeners && listeners.length
            ? listeners[listeners.length - 1].fn
            : null;
    }

    /**
     * 移除最后一个事件监听器
     * @param {string} eventName - 事件名称
     * @returns {Function|null} 被移除的监听器，不存在则返回 null
     */
    pop(eventName) {
        const listeners = this._event.get(eventName);
        if (!listeners || !listeners.length) return null;
        const removedItem = listeners.pop();
        if (!listeners.length) this._event.delete(eventName);
        return removedItem.fn;
    }

    /**
     * 删除指定的事件监听器
     * @param {string} eventName - 事件名称
     * @param {Function} listener - 要删除的监听器
     * @returns {boolean} 是否删除成功
     */
    remove(eventName, listener) {
        const listeners = this._event.get(eventName);
        if (!listeners) return false;

        /* 找到对应监听器的索引 */
        const idx = listeners.findIndex(item => item.fn === listener);
        if (idx === -1) return false;

        listeners.splice(idx, 1);
        if (!listeners.length) this._event.delete(eventName);
        return true;
    }

    /**
     * 清空指定事件的所有监听器
     * @param {string} eventName - 事件名称
     * @returns {boolean} 是否清空了事件
     */
    clear(eventName) {
        return this._event.delete(eventName);
    }

    /**
     * 清空所有事件的所有监听器
     * @returns {void}
     */
    clearAll() {
        this._event.clear();
    }

    /**
     * 获取所有事件名称
     * @returns {IterableIterator<string>} 事件名称迭代器
     */
    eventNames() {
        return this._event.keys();
    }

    /**
     * 获取指定事件的所有监听器
     * @param {string} eventName - 事件名称
     * @returns {Function[]} 监听器数组
     */
    listeners(eventName) {
        const listeners = this._event.get(eventName);
        return listeners ? listeners.map(item => item.fn) : [];
    }

    /**
     * 获取监听器数量
     * @param {string} [eventName] - 事件名称；不传则返回总数量
     * @returns {number} 数量
     */
    listenerCount(eventName) {
        if (eventName === undefined) {
            let count = 0;
            for (const arr of this._event.values()) count += arr.length;
            return count;
        }
        const listeners = this._event.get(eventName);
        return listeners ? listeners.length : 0;
    }
}

/* 挂载到全局 */
window.$bus = new EventBus();

export default window.$bus;