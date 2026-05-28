/**
 * @author 友人a丶
 * @date 2026-04-25
 * 统一拦截器模块
 */
import axios from "axios";
import useAuth from "@/stores/auth";
import {appendParamsToUrl} from "@/common/index.js";

/* 设置 */
axios.defaults.timeout = 60000;

export default function loadInterceptor() {

    /* 授权的状态 */
    const auth = useAuth();

    /* 请求拦截器 */
    axios.interceptors.request.use((req) => {
        if (auth.token) req.url = appendParamsToUrl(req.url, {token: auth.token});
        return req;
    })

    /* 响应拦截器 */
    axios.interceptors.response.use(
        (response) => {
            if (Array.isArray(response.data)) {
                response.data = {
                    code: 1,
                    errMsg: 'ok',
                    data: {
                        data: response.data,
                        current_page: 1,
                        last_page: 1,
                        per_page: response.data.length,
                        total: response.data.length,
                    },
                }
            }
            return response
        },
        (error) => Promise.reject(error),
    )
}

