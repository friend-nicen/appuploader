/**
 * @author 友人a丶
 * @date
 * 内置主题
 * */
import {useLocalStorage} from "@vueuse/core";
import {computed, inject, watch} from "vue";


/* 内置主题 */
export const themes = [
    {name: "活力橙", color: "#FF7D00"},
    {name: "浪漫红", color: "#F53F3F"},
    {name: "Bili粉", color: "#fb7299"},
    {name: "清新绿", color: "#07c160"},
    {name: "樱花粉", color: "#E84A6C"},
    {name: "晚秋红", color: "#F77234"},
    {name: "柠檬黄", color: "#EFBD48"},
    {name: "仙野绿", color: "#00B42A"},
    {name: "锌色灰", color: "#3F3F46"},
    {name: "碧涛青", color: "#14C9C9"},
    {name: "海蔚蓝", color: "#3491FA"},
    {name: "极致蓝", color: "#165DFF"},
    {name: "暗夜紫", color: "#722ED1"},
    {name: "中灰色", color: "rgb(56, 66, 82)"}
].map((item) => {
    return {
        "name": item.name,
        "primaryColor": item.color,
        "adminHeaderMenuActiveColor": item.color,
        "adminHeaderMenuColor": "#000000",
        "primaryHover": "#f7f8fa",
        "adminHeaderBgcolor": "white",
        "sidebarColor": "#303133",
        "sidebarBgColor": "#F6F9F9",
        "adminLayoutContent": "#f2f3f5",
        "sidebarHoverBgColor": "#f2f3f5",
        "cardBgColor": "#f2f3f5",
        "adminHeaderColor": "#000000d9",
        "adminHeaderBorder": "#e5e6eb",
        "sidebarHoverColor": item.color
    };
});


/**
 * 添加css变量
 * @param theme
 * @param styleId
 */
export function applyCSSVar(theme, styleId) {

    /* style */
    let cssVariables = '';

    /* 遍历所有变量 */
    for (const key in theme) {
        if (theme.hasOwnProperty(key)) {
            const value = theme[key];
            const cssVarName = `--v-${key.replace(/([a-z])([A-Z])/g, '$1-$2').toLowerCase()}`;
            cssVariables += `${cssVarName}: ${value};`;
        }
    }

    /* 检查是否已经存在具有相同id的style标签 */
    const style = document.getElementById(styleId);

    /* 如果存在，移除现有的style标签 */
    if (style) document.head.removeChild(style);

    /* 将style标签添加到head中 */
    document.head.appendChild(Object.assign(document.createElement('style'), {
        id: styleId,
        innerHTML: `:root { ${cssVariables} }`
    }));
}


/**
 * 开始动态主题监听
 */
export function dynamicTheme() {

    /* 主题 */
    const theme = useLocalStorage("theme", themes[0]);

    /* 获取APP实例 */
    const app = inject('appInstance');

    /* 监听主题变化 */
    watch(() => theme.value, () => {

        /** @type theme 主题变量 */
        const cssVar = Object.assign({
            primaryColor: '#165dff',
            linkColor: '#165dff',
        }, theme.value)

        /* 全局配置 */
        app.config.globalProperties.$theme = cssVar;

        /* 自定义 */
        applyCSSVar(cssVar, 'theme');

    }, {
        deep: true,
        immediate: true
    })

    /* 反应antd组件需要的主题色 */
    return computed(() => {
        const usedTheme = theme.value || themes[0];
        return {
            token: {
                fontSize: 14,
                colorPrimary: usedTheme.primaryColor,
                colorLink: usedTheme.primaryColor,
                colorPrimaryHover: usedTheme.primaryColor,
                colorInfo: usedTheme.primaryColor,
                colorBorderSecondary: "#0000000f"
            }
        };
    });

}