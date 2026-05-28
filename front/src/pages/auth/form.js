import { store } from '@/common';

/* 数据添加模板 */
const form_add = [
    {
        key: 'name',
        type: 'input',
        label: '名称',
        attr: {
            required: true,
            placeholder: '请输入名称 (e.g. My CI Key)'
        }
    },
    {
        key: 'issuer_id',
        type: 'input',
        label: 'Issuer ID',
        attr: {
            required: true,
            placeholder: '请输入 Issuer ID'
        }
    },
    {
        key: 'key_id',
        type: 'input',
        label: 'Key ID',
        attr: {
            required: true,
            placeholder: '请输入 Key ID'
        }
    },
    {
        key: 'private_key',
        type: 'textarea',
        label: 'Private Key',
        attr: {
            required: true,
            placeholder: '-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----',
            rows: 4
        }
    }
];

export default function initForm() {
    const need_add = store({
        name: '',
        issuer_id: '',
        key_id: '',
        private_key: ''
    });

    return {
        form_add,
        need_add
    };
}
