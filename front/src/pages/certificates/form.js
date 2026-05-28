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
        key: 'type',
        type: 'input',
        label: '类型',
        attr: {
            required: true,
            placeholder: '请输入类型 (e.g. IOS_DEVELOPMENT)'
        }
    }
];

export default function initForm() {
    const need_add = store({
        name: '',
        type: ''
    });

    return { form_add, need_add };
}
