import initTable from '@/common/init-table';

const columns = [
    { title: 'Name', dataIndex: ['attributes', 'name'] },
    { title: 'Bundle ID', dataIndex: ['attributes', 'bundleId'] },
    { title: 'SKU', dataIndex: ['attributes', 'sku'] },
    { title: 'Primary Locale', dataIndex: ['attributes', 'primaryLocale'] },
    { title: '操作', dataIndex: 'action', width: 100 }
];

export default function init() {
    const table = initTable({
        unique: 'Apps-Table',
        column: columns,
        condition: {}
    });

    return { table, columns };
}
