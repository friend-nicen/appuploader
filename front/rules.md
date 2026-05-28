---
alwaysApply: false
description: 
---
# 代码规范

1. 所有注释全部使用 /* */，的形式，放在对应代码的上方
2. components目录内的组件，框架有自动引入，不需要再次定义import
3. 所有组件全部使用vue setup进行开发

# 页面开发规范

所有页面遵循如下结构：

* index.js，导出页面组件
* table.js，初始化页面上的表格
* form.js，初始化页面上的表单
* index.vue，实现页面组件

## 页面模板如下：

```html
<a-card :bordered="false" class="card" :title="$route.meta.name">
    <template #extra>
        <!-- 这里是功能按钮 -->
    </template>
    <!-- 下面是表格 -->
</a-card>
```

## 需要用到表格组件，按照如下规范来开发：

1. 使用 src/components/v-table，这个表格组件

```vue

<v-table :init="face.list" :rowSelection="true">
<template #bodyCell="{ column, record }">
  <template v-if="column.title === '操作'">
    <v-op :source="{ column, record }"/> <!--操作-->
  </template>
  <template v-else-if="column.dataIndex === 'permission'">
    <a @click.prevent="showPermiss(record.id)">{{ record.permission }}个</a>
  </template>
</template>
</v-table>
```

2. 在页面的table.js内定义表格列，使用src/common/init-table.js 导出的方法初始化表格

```javascript
import initTable from "@/common/init-table";

/* 定义表格列 */
const columns = [
    {
        title: "路径",
        dataIndex: 'path',
        width: 180,
    },
    {
        title: "名称",
        dataIndex: 'name',
        width: 100,
    }
];

/* 初始化表格 */
const table = initTable({
    unique: 'Sys-Routes-Table',
    column: columns,
    condition: {
        name: ""
    }
});
```

3. 定义表格的筛选表单时，使用如下规范：

```javascript
/* 定义表单 */
const form = [
    {
        key: 'name',
        type: "input",
        label: "页面名称",
        attr: {
            placeholder: "请输入页面名称"
        }
    }
];
```

```vue
<!--  条件  -->
<v-batch layout="horizontal" :form="form" :data="table.condition.data"/>
```

## 需要用到表单来实现修改和增加时，按照如下规范来开发：

1. 使用 src/components/v-form，这个表格组件

```vue

<v-form
    v-model:visible="visible_add"
    :dataSource="need_add"
    :form="form_add"
    :showBorder="false"
    :init="face.add"
    :centered="true"
    :union="true"
    message="添加成功！"
    name="route_add"
    title="页面添加">
</v-form>
```

2. 在页面的form.js内定义表单：

```javascript
/* 数据添加模板 */
const form_add = [
    {
        key: 'name',
        type: "input",
        label: "页面名称",
        attr: {
            required: true,
            placeholder: "请输入页面名称",
        }
    }
];

/* 数据对象 */
const need_add = store({
    name: ""
});

/* 是否显示修改框 */
const visible_add = ref(false);

```
