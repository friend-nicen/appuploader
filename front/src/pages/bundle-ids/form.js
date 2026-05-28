import { store } from '@/common';

const platforms = [
    { label: 'iOS', value: 'IOS' },
    { label: 'macOS', value: 'MAC_OS' }
];

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
        key: 'identifier',
        type: 'input',
        label: 'Identifier',
        attr: {
            required: true,
            placeholder: '请输入 Identifier (e.g. com.example.app)'
        }
    },
    {
        key: 'platform',
        type: 'select',
        label: 'Platform',
        attr: {
            required: true,
            placeholder: '请选择平台',
            options: platforms
        }
    }
];

export default function initForm() {
    const need_add = store({
        name: '',
        identifier: '',
        platform: undefined
    });

    return { form_add, need_add };
}
