/**
 * @author 友人a丶
 * @date
 * */

import {getNum, store, useLocalStorage} from "@/common";
import {computed, provide, reactive, ref} from "vue";
import {cloneDeep} from "lodash";
import dayjs from "dayjs";
import {useRoute} from "vue-router";


/**
 * 规范表头
 * @param columns
 * @param rspv
 */
function format(columns, rspv) {

    /* 默认列配置对象 */
    const column = {
        fixed: false,
        resizable: true,
        display: true, //是否展示
        editable: false,  //可编辑
        empty: false, //允许为空
        width: 100
    };

    return columns.map((item) => {
        item = Object.assign(cloneDeep(column), item);
        if (rspv && !!item.width) {
            item.width = (item.width / 17) * window.rem
        }
        return item;
    })
}


/**
 * 递归格式化表格列配置（支持 children）
 * @param item
 * @param render
 * @return {{timestamp}|*|number}
 */
function normalize(item, render) {

    /* 自定义渲染函数 */
    if (item.dataIndex && render[item.dataIndex]) {
        item.customRender = render[item.dataIndex];
    }

    /* 排序处理 */
    if (item.sortable) {
        if (item.sortable === "number") {
            item.sorter = {
                compare: (a, b) => {
                    return getNum(a[item.dataIndex]) - getNum(b[item.dataIndex]);
                },
                multiple: 2
            };
        } else if (item.sortable === "dayjs") {
            item.sorter = {
                compare(a, b) {
                    if (dayjs(a[item.dataIndex]).isBefore(dayjs(b[item.dataIndex]))) {
                        return -1;
                    } else if (dayjs(a[item.dataIndex]).isSame(dayjs(b[item.dataIndex]))) {
                        return 0;
                    } else {
                        return 1;
                    }
                },
                multiple: 2
            };
        }
    }

    /* 时间戳格式化 */
    if (item.timestamp) {
        item.customCell = (record, rowIndex, column) => {
            const val = record[column.dataIndex];
            return {
                textContent: val ? dayjs(val * 1000).format(item.timestamp) : "-"
            };
        };
    }

    /* 递归处理子列 */
    if (item.children?.length) {
        item.children = item.children.map((i => {
            return normalize(i, render);
        }));
    }

    return item;
}

/**
 * 计算字符串唯一值
 * @param s
 * @return {string}
 * @constructor
 */
function hash(s) {
    let h = 5381 | 0, i = s.length;
    while (i) h = (h * 33 ^ s.charCodeAt(--i)) >>> 0;
    return h + "";
}

/**
 * 处理列配置和缓存
 * @param unique
 * @param column
 * @param localize
 * @return {Ref<*>}
 */
function cache(unique, column, localize) {
    /* 读取本地存储列数据 */
    const columns = useLocalStorage(unique, cloneDeep(column), localize);
    /* 生成配置哈希 */
    const value = hash(JSON.stringify(column));
    /* 读取本地缓存哈希 */
    const localHash = localStorage.getItem(unique + "-key");
    /* 配置变化则重置列 */
    if (localHash && value !== localHash) columns.value = cloneDeep(column);
    /* 更新缓存哈希 */
    localStorage.setItem(unique + "-key", value);
    /* 返回响应式列 */
    return columns;
}

/**
 * 初始化表格
 */
export default function (set) {

    /**
     * 初始化的默认配置
     * @type {{props: boolean,localize: boolean, column: null, defaultExpandAllRows: boolean, check: null|Function, filter: null, condition: null, rspv: boolean, immediate: boolean, watch: boolean, enable: *[], unique: null, effect: boolean, callback: null, id: string}}
     */
    const config = {
        id: 'table',
        unique: null, //区分本地存储的表头
        column: null, //列配置
        effect: true,
        condition: null, //查询条件
        props: false,
        methods: "POTS",
        callback: null,//数据加载完毕后的回调函数
        filter: null,//数据加载完后是否要过滤
        immediate: true, //组件初始化后立即加载数据
        watch: true, //监听条件变化，加载数据
        localize: true, //表头是否本地存储
        defaultExpandAllRows: true,//是否要展开所有行
        rspv: false, //宽度响应式，跟随rem值
        check: null, //权限校验，过滤无权限查看的列
        enable: []  // 触发排序是否重新请求
    };

    /* 重载用户的配置 */
    Object.assign(config, set);

    /* 解构配置 */
    let {
        unique,
        column,
        props,
        condition,
        callback,
        localize,
        filter,
        defaultExpandAllRows,
        check,
        rspv
    } = config


    /* 需要合并查询参数 */
    if (props && condition) Object.assign(condition, useRoute().query)


    /* 用于自定义排序的列 */
    const sort = reactive({
        left: [], plain: [], right: []
    });

    /* 处理表头 */
    column = format(column, rspv);

    /* 存在的渲染函数 */
    const render = Object.fromEntries(column.filter(o => o.customRender)
        .map(({dataIndex, customRender}) => [dataIndex, customRender]));

    /* 默认表格 */
    const defaultColumns = ref(cloneDeep(column));

    /* 读取或者更新缓存 */
    const usedColumns = cache(unique + "-used", column, localize);

    /* 触发列更新的响应式变量 */
    const trigger = ref(0);

    /* 表配置 */
    const columns = computed(() => {

        /* 开始处理 */
        let items = usedColumns.value;
        if (trigger.value < 0) return [];

        /**
         * 如果已经请求了表头
         * 这个表头是完整的包括所有数据的表头
         * 初次请求的是默认的表头（用户修改指标之后保存的表头）
         * */
        if (items) {

            const left = [];
            const plain = [];
            const right = [];
            const result = [];

            /* 一次遍历完成所有逻辑：权限 + 显示 + 格式化 + 分组 */
            for (const item of items) {

                /* 合并：权限过滤 + 不显示则跳过 */
                if ((check && !check(item)) || !item.display) {
                    continue;
                }

                /* 处理分组 */
                const col = normalize(item, render);

                /* 固定列分组 */
                if (col.fixed === "left") {
                    left.push(col);
                } else if (col.fixed === "right") {
                    right.push(col);
                } else {
                    col.fixed = false;
                    plain.push(col);
                }

                result.push(col);
            }

            /* 挂载分组 */
            Object.assign(sort, {left, plain, right});

            return result;

        }


        return [];
    });

    /* 暴露出表格初始化数据的方法 */
    const table = reactive({
        timestamp: null,
        loadData: null,
        loaded: false,
        mounted: false,
        update: () => trigger.value++,
        selectRows: {
            keys: [],
            rows: []
        },
        alter(newColumn) {
            usedColumns.value = format(newColumn, rspv);
        },
        condition: !condition ? {} : store(condition),
        callback: callback,
        filter: filter,
        watch: Object.assign({
            sorter: true,
            pagination: true,
            condition: true
        }, typeof config.watch === 'object' ? config.watch : (config.watch ? {} : {
            sorter: false,
            pagination: false,
            condition: false
        })),
        immediate: config.immediate,
        defaultExpandAllRows,
        editable: [],
        effect: config.effect,
        usedColumns,
        methods: ['get', 'post'].includes(config.methods.toLowerCase()) ? config.methods.toLowerCase() : "post",
        defaultColumns,
        enable: config.enable,
        columns,
        sort
    });

    /* 注入表格 */
    provide(config.id, table);

    return table;

}