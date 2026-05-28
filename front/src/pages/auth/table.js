import initTable from '@/common/init-table';

const columns = [
    { title: '名称', dataIndex: 'name' },
    { title: 'Key ID', dataIndex: 'key_id', class: 'font-mono' },
    { title: 'Issuer ID', dataIndex: 'issuer_id', class: 'font-mono' },
    { title: '状态', dataIndex: 'is_active', width: 100 },
    { title: '操作', dataIndex: 'action', width: 200, align: 'right' }
];

export default function init() {
    const table = initTable({
        unique: 'Auth-Keys-Table',
        column: columns,
        condition: {}
    });

    return {
        table,
        columns
    };
}
