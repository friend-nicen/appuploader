/**
 * 一些公共的方法
 * */
import dayjs, {isDayjs} from "dayjs";
import {cloneDeep, debounce} from "lodash";
import {computed, inject, isRef, provide, reactive, ref, watch} from "vue";
import {UAParser} from "ua-parser-js";

/**
 * 获取指定时间范围内的日期
 * @param start
 * @param end
 * @param format
 * @param flag
 * @returns {*[]}
 */
export function getDays(start, end, format = "YYYY-MM-DD", flag = true) {

    let days = [];
    let temp = start; //中间日期

    /* 不包含最后一天 */
    if (!flag) {
        end = end.subtract(1, 'day');
    }

    /*start大于等于end*/
    if (start.isAfter(end)) {
        return days;
    }

    /*相同那就是一天*/
    if (start.isSame(end)) {
        days.push(temp.format(format)); //默认第一天
    }

    /*遍历所有中间日期*/
    // eslint-disable-next-line
    while (1) {


        temp = temp.add(1, 'day'); //加1

        if (temp.isAfter(end)) {
            break;
        }

        days.push(temp.format(format));

    }


    return days;
}


/**
 * 延迟执行
 * @param callback
 * @param time
 */
export function delay(callback, time) {
    let timer = setTimeout(() => {
        clearTimeout(timer);
        callback();
    }, time);
}


/**
 * 禁止选择的时间范围
 * @param datetime
 * @returns {boolean}
 */
export function disabledDate(datetime) {
    let today = dayjs(); //当天的时间
    return !today.isAfter(dayjs(datetime));
}


/**
 * 返回表格行key值
 * @param record
 * @param key
 * @returns {*}
 */
export function key(record, key = "id") {
    return record[key];
}


/**
 * 获取上传文件的图片base64
 * @param file
 * @returns {Promise<unknown>}
 */
export function getBase64(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => resolve(reader.result);
        reader.onerror = error => reject(error);
    });
}

/**
 * 获取间隔某个日期的天数
 * @param day
 * @returns {number}
 */
export function diff(day) {

    let now = dayjs().format("DD");

    if (parseInt(now) > parseInt(day)) {
        let predix = dayjs(dayjs().format("YYYY-MM-") + (parseInt(day) > 10 ? day : "0" + day)).add(1, "month");
        return predix.diff(dayjs(), 'day');
    } else {
        return parseInt(day) - parseInt(now);
    }
}

/**
 * 响应式本地存储
 * @param key
 * @param value
 * @param local 是否要本地化同步
 * @returns {any|WritableComputedRef<any>}
 */
export function useLocalStorage(key, value = [], local = true) {


    /* 本地同步 */
    const item = ref(value);

    /**
     * 如果已经有值了
     * */
    if (local) {

        const primary = localStorage.getItem(key);

        if (primary) {
            item.value = JSON.parse(primary);
        } else {
            localStorage.setItem(key, JSON.stringify(value));
        }

        /**
         * 监视属性，同步整个对象到localstorage
         * 节流运行
         * */
        watch(item, debounce(() => {
            if (item.value) {
                localStorage.setItem(key, JSON.stringify(item.value));
            }
        }, 500), {
            deep: true
        });

        /**
         * 计算属性，获取和
         * */
        return computed({
            get() {
                return item.value;
            },
            set(newValue) {
                /* 保存值 */
                item.value = newValue;
                /* 如果为null，则清空本地存储 */
                if (newValue === null) {
                    localStorage.removeItem(key);
                } else {
                    localStorage.setItem(key, JSON.stringify(newValue, null, 2));
                }
            }
        })

    } else {
        return item;
    }


}


/**
 * 条件转换
 * */
export function switchForm(data, format = "YYYY-MM-DD") {

    //参数检测
    if (!data) {
        console.warn("无效参数");
        return {};
    }

    try {

        let condition;

        if (isRef(data)) {
            condition = cloneDeep(data.value);
        } else {
            condition = cloneDeep(data);
        }


        /**
         * 转换所有dayjs
         * */
        for (let i in condition) {

            if (condition[i] === "" || condition[i] === null || condition[i] === undefined) {
                delete condition[i];
                continue;
            }

            if (isDayjs(condition[i])) {
                condition[i] = condition[i].format(format);
                continue;
            }

            if (isRef(condition[i])) {
                condition[i] = condition[i].value;
            }

            if (condition[i] instanceof Array) {

                if (condition[i].length === 0) {
                    delete condition[i];
                    continue;
                }

                /* 有效数据量 */
                let count = 0;

                condition[i] = condition[i].map(item => {
                    if (isDayjs(item)) {
                        item = item.format(format);
                    }
                    if (item !== null) {
                        count++;
                    }
                    return item
                })

                if (count === 0) {
                    delete condition[i];
                }
            }
        }


        return condition;
    } catch (e) {
        console.warn(e)
    }


}


/**
 * 数据容器，并提供重置的方法
 * @param  res
 * @returns {{data: *, reset: ()=>void,$set: (object)=>void,save: (object)=>void}}
 */
export function store(res) {

    const init = reactive({
        primary: cloneDeep(res),
        data: cloneDeep(res),
        /**
         * 重置回上一个快照
         * */
        reset: () => {
            init.data = Object.assign(init.data, init.primary);
            return true;
        },
        /**
         * 替换
         * */
        $set(set) {
            init.data = Object.assign(init.data, cloneDeep(set));
            return true;
        },
        /**
         * 保存快照
         * */
        save: (snap) => {
            init.primary = cloneDeep(snap);
            return true;
        }
    })

    return init;
}


/**
 * 下拉数据搜索
 * @param input
 * @param option
 * @returns {boolean}
 */
export function filterOptions(input, option) {
    return option.label.indexOf(input) > -1;
}


/**
 * @param arr
 * @returns {*}
 * 获取数组最后一个元素
 */
export function pop(arr) {
    if (arr.length === 0) {
        return null;
    }
    return arr[arr.length - 1];
}


/**
 * 下拉数据搜索
 * @param input
 * @param treeNode
 * @returns {boolean}
 */
export function filterTreeNode(input, treeNode) {
    return treeNode.title.indexOf(input) > -1;
}


/**
 * 秒数转换
 * @param time
 * @param is_string
 * @return {(string|number)[]}
 */
export function getTime(time, is_string = false) {
    // 转换为式分秒
    let h = parseInt(time / 60 / 60 % 24)
    h = h < 10 ? '0' + h : h
    let m = parseInt(time / 60 % 60)
    m = m < 10 ? '0' + m : m
    let s = parseInt(time % 60)
    s = s < 10 ? '0' + s : s
    // 作为返回值返回

    if (is_string) {
        return `${h}时${m}分${s}秒`
    } else {
        return [h, m, s]
    }


}


/**
 * 元素toggle
 * @param arr
 * @param element
 * @return {*}
 */
export function arrayToggle(arr, element) {
    /* 查找元素在数组中的索引 */
    const index = arr.indexOf(element);
    /* 如果元素存在，则移除 */
    if (index > -1) {
        arr.splice(index, 1);
    } else {
        // 如果元素不存在，则添加
        arr.push(element);
    }
    // 返回更新后的数组
    return arr;
}


/**
 * 批量inject
 */
export function injects(keys) {
    const obj = Object.create(null);
    for (let i of keys) {
        obj[i] = inject(i, null);
    }
    return obj;
}

/**
 * 批量provide
 */
export function provides(obj) {
    for (let i in obj) {
        provide(i, obj[i])
    }
}

/**
 * 休眠
 * @param time
 * @returns {Promise<unknown>}
 */
export function sleep(time) {
    return new Promise(resolve => {
        let timer = setTimeout(() => {
            clearTimeout(timer);
            resolve();
        }, time);
    })
}


/**
 * 获取get参数
 * @return {{}}
 */
export function getQueryParams() {

    /* 获取完整的URL */
    const url = window.location.href;

    /* 定义一个对象来存储参数 */
    const params = {};

    /* 定义一个函数来解析参数字符串 */
    function parseParams(paramStr) {
        if (!paramStr) return;
        paramStr.split("&").forEach((param) => {
            const [key, value] = param.split("=");
            if (key) {
                params[decodeURIComponent(key)] = decodeURIComponent(value || "");
            }
        });
    }


    const firstQuestion = url.indexOf("?");
    const firstHash = url.indexOf("#");

    let queryString = "";

    if (firstQuestion !== -1) {

        /* 如果存在 #，则截取 ? 到 # 之间的部分 */
        if (firstHash !== -1 && firstHash > firstQuestion) {
            queryString = url.substring(firstQuestion + 1, firstHash);
        } else {
            /*  如果没有 #，则截取 ? 到 URL 结尾的部分 */
            queryString = url.substring(firstQuestion + 1);
        }
        /* 将第一个 ? 后面的数据中的 ? 替换为 &（如果需要） */
        const modifiedQueryString = queryString.replace(/\?/g, "&");
        parseParams(modifiedQueryString);
    }


    /* 解析URL中的#后面的GET参数 */
    const hashString = url.split("#")[1];
    if (hashString && hashString.includes("?")) {
        const hashParams = hashString.split("?")[1];
        parseParams(hashParams);
    }

    return params;
}


/**
 * 获取UA
 * @param ua
 * @return {string}
 */
export function getUa(ua) {
    if (!ua) return '-';
    try {
        const {device, os, browser} = UAParser(ua);
        const {vendor, model} = device;
        if (vendor) {
            return `${vendor} ${model || ''}`.trim();
        } else if (os.name) {
            return `${os.name || ''} ${browser.name || ''}`.trim();
        } else {
            return ua;
        }
    } catch (e) {
         
        void e;
        return ua;
    }
}


/**
 * 下载Blob对象
 * @param blob
 * @param fileName
 */
export function downloadBlob(blob, fileName = null) {
    const url = URL.createObjectURL(blob);
    const download = document.createElement('a');
    download.href = url;
    download.download = fileName ? fileName : (new Date()).getTime() + '.png';
    download.click();
    URL.revokeObjectURL(url);
}


/**
 * 按照对象值去重
 * @param arr
 * @param key
 * @return {*}
 */
export function uniqueByKey(arr, key) {
    const seen = new Set(); // 用于存储已经出现过的 key 值
    return arr.filter(item => {
        const keyValue = item[key];
        if (!seen.has(keyValue)) {
            seen.add(keyValue); // 如果 key 值未出现过，添加到 Set 中
            return true; // 保留当前项
        }
        return false; // 如果 key 值已出现过，过滤掉当前项
    });
}


/**
 *  提取值中的有效数字（支持小数），无有效数字返回 0
 **/
export function getNum(v) {
    if (v == null || v === '') return 0;
    const numMatch = String(v).match(/\d+(\.\d+)?/);
    return numMatch ? Number(numMatch[0]) || 0 : 0;
}


/**
 * 随机打乱数组（Fisher-Yates 洗牌算法，公平无偏倚）
 * @param arr
 * @return {*[]}
 */
export function shuffleArray(arr) {
    // 先深拷贝原数组，避免修改原始数组（可选，根据需求决定是否保留）
    const newArr = [...arr];
    for (let i = newArr.length - 1; i > 0; i--) {
        /* 生成 0 到 i 之间的随机索引 */
        const j = Math.floor(Math.random() * (i + 1));
        /* 交换当前元素和随机索引元素 */
        [newArr[i], newArr[j]] = [newArr[j], newArr[i]];
    }
    return newArr;
}


/**
 * 计算总和
 * @param data
 * @param key
 * @returns {*}
 */
export function getSum(data, key) {
    return data.reduce((total, item) => {
        return total + (Number(item[key]) || 0)
    }, 0)
}

/**
 * 获取纯文本
 * @param html
 * @return {string|string}
 */
export function stripHtml(html) {
    const doc = new DOMParser().parseFromString(html, 'text/html');
    return doc.body.textContent || '';
}



/**
 * 增加参数
 * @param url
 * @param params
 * @return {string}
 */
export function appendParamsToUrl(url, params) {

    /* 检查链接是否已经有查询参数 */
    const hasParams = url.includes('?');

    /* 遍历参数对象，拼接为查询字符串 */
    const queryString = Object.keys(params)
        .map(key => `${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`)
        .join('&');

    /* 根据链接是否有查询参数，决定如何拼接 */
    if (hasParams) {
        /* 果已经有参数，添加新的参数时用 '&' 连接 */
        return `${url}&${queryString}`;
    } else {
        /* 如果没有参数，添加新的参数时用 '?' 连接 */
        return `${url}?${queryString}`;
    }
}


