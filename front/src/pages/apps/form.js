import { store } from '@/common';

const form_add = [
    {
        key: 'name',
        type: 'input',
        label: '名称',
        attr: {
            required: true,
            placeholder: '请输入名称'
        }
    },
    {
        key: 'bundleId',
        type: 'input',
        label: 'Bundle ID',
        attr: {
            required: true,
            placeholder: '请输入 Bundle ID'
        }
    }
];

export default function initForm() {
    const need_add = store({
        name: '',
        bundleId: ''
    });

    return { form_add, need_add };
}
