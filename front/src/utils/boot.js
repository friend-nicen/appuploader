/**
 * @author 友人a丶
 * @date 2026-04-25
 *
 * 引导系统初始化
 * 初始化全局响应拦截器
 * 初始化路由守卫
 * 初始化用户登录（预留）
 */

import loadGuard from '@/router/guard/load-guard'
import loadInterceptor from '@/service/interceptor'
import timezone from 'dayjs/plugin/timezone'
import utc from 'dayjs/plugin/utc'
import 'nprogress/nprogress.css'
import dayjs from 'dayjs'

export default async function boot() {

  /* 加载拦截器 */
  loadInterceptor();

  /* 加载路由守卫 */
  loadGuard();

  /* 加载 dayjs 插件 */
  dayjs.extend(utc);
  dayjs.extend(timezone);
}

