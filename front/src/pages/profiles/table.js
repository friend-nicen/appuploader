import initTable from '@/common/init-table';

const columns = [
    { title: 'ID', dataIndex: 'id', width: 120 },
    { title: '名称', dataIndex: ['attributes', 'name'] },
    { title: '类型', dataIndex: ['attributes', 'profileType'] },
    { title: 'Platform', dataIndex: ['attributes', 'platform'] },
    { title: '状态', dataIndex: ['attributes', 'profileState'] },
    { title: '过期时间', dataIndex: ['attributes', 'expirationDate'] },
    { title: '操作', dataIndex: 'action', width: 80, empty: true }
];

export default function init() {
    const table = initTable({
        unique: 'Profiles-Table',
        column: columns,
        condition: {}
    });

    return { table, columns };
}
