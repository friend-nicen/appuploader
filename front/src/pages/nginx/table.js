import initTable from '@/common/init-table'
import {h, reactive, resolveComponent} from 'vue'
import api from '@/service/api'

/* 日志查询页面初始化 */
export default function init() {

    /* 状态码颜色 */
    const statusColor = (status) => {
        const s = Number(status)
        if (s >= 200 && s < 300) return 'green'
        if (s >= 300 && s < 400) return 'blue'
        if (s >= 400 && s < 500) return 'orange'
        if (s >= 500) return 'red'
        return 'default'
    }


    /* 定义表格列（默认展示所有字段） */
    const columns = [
        {
            title: '时间',
            dataIndex: 'time',
            width: 170,
            empty: true,
        },
        {
            title: '状态码',
            dataIndex: 'status',
            width: 90,
            sortable: 'number',
            customRender: ({text}) => {
                return h(
                    resolveComponent('a-tag'),
                    {color: statusColor(text)},
                    {default: () => String(text)},
                )
            },
        },
        {title: '域名', dataIndex: 'domain', width: 160},
        {title: '请求方法', dataIndex: 'method', width: 90},
        {title: '请求地址', dataIndex: 'full_url', width: 210},
        {title: '用户IP', dataIndex: 'remote_addr', width: 140},
        {title: 'UA', dataIndex: 'ua_parsed', width: 125},
        {title: '反代理', dataIndex: 'upstream_addr', width: 160},
        {title: 'Referer', dataIndex: 'referer', width: 210},
        {
            title: '请求耗时',
            dataIndex: 'req_time',
            width: 120,
            sortable: 'number',
            empty: true,
        },
        {
            title: '接口耗时',
            dataIndex: 'up_time',
            width: 120,
            sortable: 'number',
            empty: true,
        },

    ]

    /* 顶部统计（仅展示） */
    const summary = reactive({
        window: '5m',
        updatedAt: '',
        total: 0,
        ipCount: 0,
        status: {'2xx': 0, '3xx': 0, '4xx': 0, '5xx': 0},
        topRoutes: [],
    })


    /* 常用筛选表单（首屏展示） */
    const form = [
        {
            key: 'window',
            type: 'select',
            label: '时间窗口',
            attr: {
                placeholder: '选择相对时间窗口',
                options: [
                    {label: '5m', value: '5m'},
                    {label: '15m', value: '15m'},
                    {label: '1h', value: '1h'},
                    {label: '6h', value: '6h'},
                    {label: '24h', value: '24h'},
                ],
            },
        },
        {
            key: 'domain',
            type: 'input',
            label: '域名',
            attr: {placeholder: '域名模糊匹配'},
        },
        {
            key: 'route',
            type: 'input',
            label: '路由/URL',
            attr: {placeholder: '路由/URL 模糊匹配'},
        },
        {
            key: 'status',
            type: 'input-number',
            label: '状态码',
            attr: {placeholder: '200/404/500...'},
        },
        {
            key: 'ip',
            type: 'input',
            label: '访客IP',
            attr: {placeholder: '访客 IP 模糊匹配'},
        },
        {
            key: 'datetime',
            type: 'range-picker',
            label: '访问时间',
            attr: {
                allowEmpty: [true, true],
                showTime: {format: 'HH:mm:ss'},
            },
        },
        {
            key: 'referer',
            type: 'input',
            label: '来源地址',
            attr: {placeholder: 'referer 模糊匹配'},
        },
        {
            key: 'req',
            type: 'input-range',
            label: '总耗时(s)'
        },
        {
            key: 'up',
            type: 'input-range',
            label: '后端耗时(s)'
        }
    ];


    /* 初始化表格对象 */
    const table = initTable({
        unique: 'Nginx-Log-Monitor-Table',
        column: columns,
        props: true,
        watch: {
            pagination: false
        },
        condition: {
            window: '15m',
            datetime: [],
            domain: null,
            route: null,
            status: null,
            ip: null,
            referer: null,
            req: [],
            up: []
        },
        callback(res) {

            /* 统计信息 */
            summary.total = res.data.data.length
            summary.window = table.condition.data.window || summary.window
            summary.updatedAt = new Date().toLocaleString()

            /* 状态数据 */
            const statusCount = {'2xx': 0, '3xx': 0, '4xx': 0, '5xx': 0}
            const ipSet = new Set()
            const routeMap = new Map()

            /* 统计数据 */
            res.data.data.forEach((r) => {
                const ip = r.remote_addr || ''
                if (ip) ipSet.add(ip)
                const s = Number(r.status || 0)
                if (s >= 200 && s < 300) statusCount['2xx'] += 1
                else if (s >= 300 && s < 400) statusCount['3xx'] += 1
                else if (s >= 400 && s < 500) statusCount['4xx'] += 1
                else if (s >= 500 && s < 600) statusCount['5xx'] += 1
                const route = r.route || r.full_url || '-'
                const k = String(route)
                routeMap.set(k, (routeMap.get(k) || 0) + 1)
            })

            /* 统计数据 */
            summary.ipCount = ipSet.size
            summary.status = statusCount
            summary.topRoutes = Array.from(routeMap.entries())
                .sort((a, b) => b[1] - a[1])
                .slice(0, 3)
                .map(([route, count]) => ({key: route, route, count}))
        }
    })

    return {
        table,
        face: {
            metrics: api.metrics,
            health: api.health
        },
        summary,
        columns,
        form
    }
}
