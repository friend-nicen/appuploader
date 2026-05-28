import { store } from '@/common';

/* Profile Type 选项列表 */
const profileTypes = [
    { label: 'iOS App Development', value: 'IOS_APP_DEVELOPMENT' },
    { label: 'iOS App Store', value: 'IOS_APP_STORE' },
    { label: 'iOS App AdHoc', value: 'IOS_APP_ADHOC' },
    { label: 'iOS App InHouse', value: 'IOS_APP_INHOUSE' },
    { label: 'macOS App Development', value: 'MAC_APP_DEVELOPMENT' },
    { label: 'macOS App Store', value: 'MAC_APP_STORE' },
    { label: 'macOS App Direct', value: 'MAC_APP_DIRECT' },
    { label: 'Mac Catalyst Development', value: 'MAC_CATALYST_APP_DEVELOPMENT' },
    { label: 'Mac Catalyst Store', value: 'MAC_CATALYST_APP_STORE' },
    { label: 'Mac Catalyst Direct', value: 'MAC_CATALYST_APP_DIRECT' },
];

/* Bundle ID 字段索引，供外部动态设置 options */
const BUNDLE_ID_INDEX = 2;
/* 证书字段索引 */
const CERT_INDEX = 3;
/* 设备字段索引 */
const DEVICE_INDEX = 4;

const form_add = [
    {
        key: 'name',
        type: 'input',
        label: '名称',
        attr: {
            required: true,
            placeholder: '请输入描述文件名称'
        }
    },
    {
        key: 'profileType',
        type: 'select',
        label: '类型',
        attr: {
            required: true,
            placeholder: '请选择描述文件类型',
            options: profileTypes
        }
    },
    {
        key: 'bundleId',
        type: 'select',
        label: 'Bundle ID',
        attr: {
            required: true,
            placeholder: '请选择 Bundle ID',
            options: []
        }
    },
    {
        key: 'certIds',
        type: 'multi-select',
        label: '证书',
        attr: {
            required: true,
            placeholder: '请选择证书',
            options: []
        }
    },
    {
        key: 'deviceIds',
        type: 'multi-select',
        label: '设备',
        attr: {
            placeholder: '请选择设备（可选）',
            options: []
        }
    }
];

export default function initForm() {
    const need_add = store({
        name: '',
        profileType: undefined,
        bundleId: undefined,
        certIds: [],
        deviceIds: []
    });

    return { form_add, need_add };
}

export { BUNDLE_ID_INDEX, CERT_INDEX, DEVICE_INDEX };
