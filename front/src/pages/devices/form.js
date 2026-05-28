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
        key: 'udid',
        type: 'input',
        label: 'UDID',
        attr: {
            required: true,
            placeholder: '请输入设备 UDID'
        }
    }
];

export default function initForm() {
    const need_add = store({
        name: '',
        udid: ''
    });

    return { form_add, need_add };
}
