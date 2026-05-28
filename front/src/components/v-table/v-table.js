/* eslint-disable */
/**
 * @author 友人a丶
 * @date
 * */

import {computed, inject, onActivated, onDeactivated, onMounted, reactive, ref, toRef, watch} from "vue";
import {switchForm} from "@/common";
import {cloneDeep, debounce} from 'lodash';
import {getNode} from "@/common/tree";
import load from "@/common/load";
import axios from "axios";
import dayjs from "dayjs";
import qs from "qs";

export default function (props) {

    /* 依赖 */
    const table = inject(props.tid);
    const expandedRowKeys = ref([]);
    const loaded = ref(false);

    /* 显示 */
    let show = true;

    /* 激活 */
    onActivated(() => {
        table.loadData = debounce(loadData, 300);
        table.dataSource = dataSource;
        show = true;
    });

    /* 隐藏 */
    onDeactivated(() => {
        show = false;
    })


    /**
     * 返回表格行key值
     * @param record
     * @returns {*}
     */
    function key(record) {
        return record[props.rowsKey];
    }

    /**
     * 表数据
     * */
    const dataSource = reactive({
        data: {
            data: []
        }
    })

    /* 挂载到对象 */
    table.dataSource = dataSource;

    /**
     * 如果指定了数据源
     * */
    if (!!props.dataSource) {
        dataSource.data.data = computed(() => {

            /* 最终输出的数据 */
            let data = [];

            /**
             * 筛选条件
             * */

            let condition = switchForm(table.condition.hasOwnProperty('data') ? table.condition.data : {});

            /**
             * 筛选遍历
             * */
            props.dataSource.forEach(item => {

                let result = true; //筛选结果

                /* 条件过滤 */
                for (let i in condition) {
                    if (item[i].indexOf(condition[i]) === -1) {
                        result = false;
                    }
                }

                if (result) data.push(item);
            })

            return data;

        })
    }


    /**
     * 分页
     * */
    const pagination = reactive(Object.assign({
        current: 1,
        pageSize: 30,
        pageSizeOptions: ['30', '50', '100', '300', '500'],
        showSizeChanger: true,
        showQuickJumper: true,
        lastPage: 0,
        total: 0
    }, props.pagination))

    /* 挂载 */
    table.pagination = pagination;

    /**
     * 初始化表头
     * */
    const column = toRef(table, "columns");

    /* 加载效果 */
    const loadingTable = ref(false);

    /* 可编辑的键值 */
    const editKey = Object.create(null);


    /* 计算 */
    column.value.filter(item => {
        return item.editable;
    }).forEach(item => {
        editKey[item.dataIndex] = false;
    });

    /**
     * 初始化数据
     * @param paginate
     * @param filters
     * @param sorter
     */

    const loadData = function (paginate = null, filters, sorter) {

        /* 隐藏时不加载 */
        if (!show) return;

        /* 排序重新加载 */
        if (sorter && Object.keys(sorter).length &&
            table.enable.indexOf('sort') === -1 &&
            !(!paginate || (paginate && JSON.stringify(paginate) !== JSON.stringify(pagination)))) {
            return
        }


        /**
         * 翻页的页号
         * */
        let page = pagination.current;

        /* 排序 */
        sorter = !!sorter ? (Array.isArray(sorter) ? {
            order: {
                field: sorter[0].field,
                by: sorter[0].order
            }
        } : {
            order: {
                field: sorter.field,
                by: sorter.order
            }
        }) : {}; //排序

        /* 分页计算 */
        if (paginate) {
            if (!!paginate.current) {
                page = paginate.current;
                pagination.current = paginate.current;
                pagination.pageSize = paginate.pageSize;
            }
            /* 不监听分页 */
            if (paginate && !table.watch.pagination) return;
        } else {
            pagination.current = 1;
        }

        /* 如果传递了静态数据 */
        if (!!props.dataSource) return;

        /* 加载效果 */
        if (table.effect) loadingTable.value = true;

        /* 合并请求参数 */
        const form = Object.assign({
            pageSize: pagination.pageSize,
            page: page
        }, sorter, switchForm(table.condition.hasOwnProperty('data') ? table.condition.data : {}));

        /**
         开始请求
         合并请求条件
         */
        let requestPromise;
        if (typeof props.init === 'function') {
            requestPromise = Promise.resolve(props.init(form)).then(res => {
                let parsed = res;
                if (typeof res === 'string') {
                    try {
                        parsed = JSON.parse(res);
                    } catch (e) {
                    }
                }
                let dataObj = parsed.data || parsed;
                let responseData = {
                    code: parsed.code !== undefined ? parsed.code : 1,
                    data: Array.isArray(dataObj) ? {data: dataObj, total: dataObj.length} : dataObj,
                    errMsg: parsed.errMsg || ''
                };
                return {data: responseData};
            });
        } else {
            requestPromise = table.methods === 'post' ?
                axios.post(props.init, form) :
                axios.get(`${props.init}?` + qs.stringify(form));
        }

        return requestPromise.then((res) => {

            console.log(res);

            /* 判断请求结果 */
            if (res.data.code) {

                /* 是否需要过滤数据 */
                if (!!table.filter) {
                    dataSource.data = table.filter(res.data.data.data);
                } else {
                    dataSource.data = res.data.data;
                }

                /* 是否需要展开行 */
                if (!!table.defaultExpandAllRows) {
                    expandedRowKeys.value = getNode(dataSource.data.data, 'id');
                }

                /* 更新分页信息 */
                pagination.total = res.data.data.total;

                /* 实时更新 */
                if (table.watch.pagination) {
                    pagination.current = res.data.data.current_page;
                    pagination.lastPage = res.data.data.last_page;
                    pagination.pageSize = res.data.data.per_page;
                }

                /* 请求结束 */
                if (!!table.callback) {
                    table.callback(res.data, pagination);
                }

                /* 可编辑的列和索引 */
                table.editable = dataSource.data.data.map(() => {
                    return cloneDeep(editKey);
                });

            } else {
                /* 弹出错误原因 */
                load.error(res.data.errMsg);

                /* 错误回调 */
                if (!!table.error) {
                    table.error(res.data, pagination);
                }

            }
        }).catch((e) => {
            /* 弹出错误原因 */
            load.error(e.message);
        }).finally(() => {
            /* 关闭加载效果 */
            loadingTable.value = false;
            loaded.value = true;
            table.loaded = true;
            /* 加载触发的时间戳 */
            table.timestamp = dayjs().unix();
            /* 关闭加载信息 */
            load.loaded();
        });
    }


    /**
     * 列宽调整
     * */
    const handleResizeColumn = (w, col) => {
        col.width = w;
    };


    /* 是否需要加载数据 */
    if (table.immediate) {
        onMounted(loadData);
    } else {
        onMounted(() => table.mounted = true);
    }

    /* 注入方法 */
    table.loadData = debounce(loadData, 300);


    /**
     * 表格是否具有条件
     * 监控条件变化
     * */
    if (!!table.condition.data) {
        watch(table.condition.data, () => {
            /* 如果停止监听 */
            if (!table.watch.condition) return;
            /* 筛选条件变化。重置页号 */
            pagination.current = 1;
            table.loadData();
        }, {
            deep: true
        })
    }

    /**
     * 列选择触发
     * */
    if (props.rowSelection) {

        let rowsAll = []; //所有被选中的列

        var selectConfig = {
            onChange(keys, rows) {
                table.selectRows.keys = keys;
                rowsAll = rows;
            },
            selectedRowKeys: computed(() => {
                return table.selectRows.keys;
            })
        }

        /**
         * keys被删除时，同步rows的状态
         * */
        table.selectRows.rows = computed(() => {
            let keys = table.selectRows.keys;
            return rowsAll.filter(item => {
                return keys.indexOf(item[props.rowsKey]) > -1
            })
        });
    }


    /* 表格布局 */
    const layout = computed(() => {
        if (dataSource.data.data.length === 0) {
            return "auto";
        } else {
            return "fixed";
        }
    });


    return {
        column,
        data: dataSource,
        key,
        loadData: table.loadData,
        loadingTable,
        paginate: pagination,
        handleResizeColumn,
        selectConfig,
        expandedRowKeys,
        layout,
        loaded
    }
}
