import { createApp } from 'vue';
import { createPinia } from 'pinia';
import Antd from 'ant-design-vue';
import App from './App.vue';

import { router } from './router';
import goto____bootstrap from './utils/boot';
import register____global from './utils/register';

/* ant-design-vue v4 */
import 'ant-design-vue/dist/reset.css'

/* 初始化vue */
const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(Antd);

goto____bootstrap();
register____global(app);

app.mount('#app')
