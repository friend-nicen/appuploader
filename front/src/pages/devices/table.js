import initTable from '@/common/init-table';

const columns = [
    { title: '名称', dataIndex: ['attributes', 'name'] },
    { title: '设备类 (Class)', dataIndex: ['attributes', 'deviceClass'] },
    { title: 'Platform', dataIndex: ['attributes', 'platform'] },
    { title: 'UDID', dataIndex: ['attributes', 'udid'] },
    { title: '状态', dataIndex: ['attributes', 'status'] },
    { title: '操作', dataIndex: 'action', width: 100 }
];

export default function init() {
    const table = initTable({
        unique: 'Devices-Table',
        column: columns,
        condition: {}
    });

    return { table, columns };
}
