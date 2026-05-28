/**
 * @友人a丶
 * 弹出层的简单封装
 * */

import {Button, message, Modal, notification} from "ant-design-vue";
import {delay} from "lodash";
import {h} from "vue"

let hide = 0; //弹出的个数
let key = "loading"; //弹窗的key值
let close = null; //关闭的回调

export default {

    loading(text = '加载中...') {
        close = message.loading({
            content: text,
            key: key,
            duration: 0
        });
        hide++;
    },
    loaded(flag = false) {
        /* 强制关闭 */
        if (flag) hide = 0;
        /* 是否可以关闭 */
        if (hide > 0) {
            hide--;
            if (hide === 0 && !!close) close();
        } else {
            if (close) close();
        }
    },
    error(text = '加载异常') {
        console.log((new Error()).stack)
        message.error(text, 5);
    },
    success(text = 'ok!') {
        message.success(text, 5);
    },
    info(text, config) {

        return Modal.info(Object.assign({
            title: '提示',
            content: h('div', {style: 'font-size:15px;line-height:1.8;', innerHTML: text}),
            maskClosable: false
        }, config));

    },
    warning(text, config) {
        return Modal.warning(Object.assign({
            title: '提示',
            content: h('div', {style: 'font-size:15px;line-height:1.8;', innerHTML: text}),
            maskClosable: false
        }, config));
    },
    succeed(text, config) {

        return Modal.success(Object.assign({
            title: '提示',
            content: h('div', {style: 'font-size:15px;line-height:1.8;', innerHTML: text}),
            maskClosable: false
        }, config));

    },
    confirm(text, callback = null, cancel = null, config) {
        return Modal.confirm(Object.assign({
            title: '提示',
            content: text,
            maskClosable: false,
            onOk: (close) => {
                close(); //关闭
                if (callback) {
                    callback()
                }
            },
            onCancel() {
                if (cancel) {
                    cancel()
                }
            }
        }, config))
    },
    /* 右侧通知卡片*/
    notify(config) {

        /* 配置 */
        const param = Object.assign({
            type: 'success',
            title: '温馨提示',
            message: "",
            duration: 5,
            onOk: null,
            buttonText: "确认"
        }, typeof config === 'string' ? {
            message: config
        } : config);

        /* 唯一Key */
        const key = `open${Date.now()}`;

        /* 触发通知 */
        notification[param.type](Object.assign({
            key: key,
            duration: 0,
            message: param.title,
            description: h('div', {style: 'font-size:15px;line-height:1.8;', innerHTML: param.message}),
            btn: () =>
                h(
                    Button,
                    {
                        type: 'primary',
                        size: 'small',
                        onClick: () => {
                            notification.close(key);
                            if (param.onOk) param.onOk()
                        },
                    },
                    {default: () => param.buttonText},
                ),
        }, param.option ? param.option : {}));

        /* 定时关闭 */
        if (param.duration) {
            delay(() => notification.close(key), param.duration * 1000);
        }

    }
}
