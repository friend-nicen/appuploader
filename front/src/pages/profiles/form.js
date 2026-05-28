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
        key: 'profileType',
        type: 'input',
        label: 'Profile Type',
        attr: {
            required: true,
            placeholder: '请输入 Profile Type (e.g. IOS_APP_DEVELOPMENT)'
        }
    },
    {
        key: 'bundleId',
        type: 'input',
        label: 'Bundle ID',
        attr: {
            required: true,
            placeholder: '请输入 Bundle ID 的 ID'
        }
    },
    {
        key: 'certId',
        type: 'input',
        label: 'Certificate ID',
        attr: {
            required: true,
            placeholder: '请输入 Certificate 的 ID'
        }
    }
];

export default function initForm() {
    const need_add = store({
        name: '',
        profileType: '',
        bundleId: '',
        certId: ''
    });

    return { form_add, need_add };
}
