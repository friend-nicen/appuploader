import initTable from '@/common/init-table';

const columns = [
    { title: '名称', dataIndex: ['attributes', 'name'] },
    { title: 'Identifier', dataIndex: ['attributes', 'identifier'] },
    { title: 'Platform', dataIndex: ['attributes', 'platform'] },
    { title: 'ID', dataIndex: 'id' },
    { title: '操作', dataIndex: 'action', width: 100 }
];

export default function init() {
    const table = initTable({
        unique: 'BundleIds-Table',
        column: columns,
        condition: {}
    });

    return { table, columns };
}
