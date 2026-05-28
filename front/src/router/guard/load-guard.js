/**
 * @author 友人a丶
 * @date 2022-07-11
 * */
import load_sys from "@/service/load-sys.js";
import NProgress from 'nprogress'
import sys from "@/stores/sys.js";
import router from "@/router/index.js";

/* 进度条 */
NProgress.configure({showSpinner: false})

/**
 * 加载路由守卫
 * */
export default function () {

    /**
     * 加载进度条
     * */
    router.beforeEach((to, from, next) => {
        if (!NProgress.isStarted()) NProgress.start()
        next()
    });

    /**
     * 判断系统是否初始化
     * */
    router.beforeEach(async (to, from, next) => {
        /* 系统尚未初始化 */
        if (!sys.loaded) await load_sys();
        next(); //下一个
    });


    /**
     * 切换页面标题
     * */
    router.beforeEach((to, from, next) => {
        document.title = to.meta.name
        next();
    })


    /**
     * 结束进度条
     * */
    router.afterEach(() => {
        NProgress.done();
    });


}

