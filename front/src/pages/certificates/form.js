import { store } from '@/common';

/* 证书类型选项列表 */
const certificateTypes = [
    { label: 'iOS Development', value: 'IOS_DEVELOPMENT' },
    { label: 'iOS Distribution', value: 'IOS_DISTRIBUTION' },
    { label: 'macOS App Development', value: 'MAC_APP_DEVELOPMENT' },
    { label: 'macOS App Distribution', value: 'MAC_APP_DISTRIBUTION' },
    { label: 'macOS Installer Distribution', value: 'MAC_INSTALLER_DISTRIBUTION' },
    { label: 'Mac Catalyst Development', value: 'MAC_CATALYST_DEVELOPMENT' },
    { label: 'Mac Catalyst Distribution', value: 'MAC_CATALYST_DISTRIBUTION' },
];

const form_add = [
    {
        key: 'name',
        type: 'input',
        label: '名称',
        attr: {
            required: true,
            placeholder: '请输入证书名称'
        }
    },
    {
        key: 'type',
        type: 'select',
        label: '证书类型',
        attr: {
            required: true,
            placeholder: '请选择证书类型',
            options: certificateTypes
        }
    },
    {
        key: 'password',
        type: 'input-password',
        label: '密码',
        attr: {
            required: true,
            placeholder: '请输入导出证书密码（至少6位）',
            rules: [
                { min: 6, message: '密码至少 6 位' }
            ]
        }
    }
];

export default function initForm() {
    const need_add = store({
        name: '',
        type: undefined,
        password: ''
    });

    return { form_add, need_add };
}