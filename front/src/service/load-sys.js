/**
 * @author 友人a丶
 * @date 2022-07-11
 * 初始化用户信息
 * */
import {getQueryParams} from "@/common/index.js";
import load from "@/common/load";
import sys from "@/stores/sys";
import auth from "@/stores/auth.js";


export default async function () {


    load.loading("初始化....");

    const query = getQueryParams();

    auth().token = query?.token;

    /* 加载完毕 */
    load.loaded();

    /* 标记系统初始化完毕 */
    sys.loaded = true;


    return true;

}