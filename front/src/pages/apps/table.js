import initTable from '@/common/init-table';

const columns = [
    { title: 'ID', dataIndex: 'id', width: 120 },
    { title: 'Name', dataIndex: ['attributes', 'name'] },
    { title: 'Bundle ID', dataIndex: ['attributes', 'bundleId'] },
    { title: 'SKU', dataIndex: ['attributes', 'sku'] },
    { title: 'Primary Locale', dataIndex: ['attributes', 'primaryLocale'] },
    { title: 'Link', dataIndex: ['links', 'self'], width: 300, ellipsis: true }
];

export default function init() {
    const table = initTable({
        unique: 'Apps-Table',
        column: columns,
        condition: {}
    });

    return { table, columns };
}
