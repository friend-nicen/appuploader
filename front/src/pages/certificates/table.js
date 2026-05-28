import initTable from '@/common/init-table';

const columns = [
    { title: 'ID', dataIndex: 'id', width: 120 },
    { title: '名称', dataIndex: ['attributes', 'name'] },
    { title: '类型', dataIndex: ['attributes', 'certificateType'] },
    { title: '序列号', dataIndex: ['attributes', 'serialNumber'] },
    { title: '过期时间', dataIndex: ['attributes', 'expirationDate'] },
    { title: '操作', dataIndex: 'action', width: 160, empty: true }
];

export default function init() {
    const table = initTable({
        unique: 'Certificates-Table',
        column: columns,
        condition: {}
    });

    return { table, columns };
}
