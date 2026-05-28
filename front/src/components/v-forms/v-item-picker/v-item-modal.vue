<template>
  <a-modal
      :open="visible"
      :maskClosable="false"
      title="图标选则"
      :width="800"
      class="antd-reset-modal"
      destroy-on-close
      @cancel="close">


    <template #footer>
      <div class="footer">
          <span class="info">
            <span style="color: red">*</span>
              <span class="info">已填充{{ dataSource.length }}组数据</span>
          </span>
        <a-space>
          <v-button @click="close">取消</v-button>
          <v-button type="primary" @click="add_data">添加</v-button>
          <v-button type="primary" @click="submit">确定</v-button>
        </a-space>
      </div>
    </template>

    <div class="container">
      <v-table :pagination="false" :sticky="{offsetHeader:0}" :dataSource="dataSource">
        <template #bodyCell="{ column,index,record }">

          <template v-if="column.title === '操作'">
            <a @click.prevent="del_col(index)">删除</a>
          </template>

          <template v-else-if="column.title === '序号'">
            {{ record.id }}
          </template>

          <template v-else>
            <a-input required v-model:value="record[column.dataIndex]" :placeholder="`请输入${column.title}`"/>
          </template>
        </template>
      </v-table>
    </div>
  </a-modal>
</template>

<script setup>
import {ref, inject, watch} from "vue";
import initTable from '@/common/init-table'
import load from "@/common/load";


/* 属性 */
const props = defineProps({
  'visible': {
    required: true
  },
  'unique': {
    required: true
  },
  'columns': {
    required: true
  }
});


/* 事件注入 */
const emit = defineEmits(['update:visible']);
const data = inject("item-picker-data"); //依赖注入
const update = inject("update"); //更新值


/**
 * 初始化表格
 * */
initTable({
  unique: props.unique,  //区分本地存储的表头
  column: props.columns, //列配置
  localize: false
});

/**
 * 列索引
 * */
const cols = props.columns.filter(item => {
  if (!item.dataIndex || item.dataIndex === "id") {
    return false;
  } else {
    return true;
  }
}).map(item => {
  return item.dataIndex;
})


/* 数据转换 */
const dataSource = ref([]); //数据源

/* 同步数据 */
watch(() => props.visible, () => {

  if (data.value.length > 0 && props.visible) {
    dataSource.value = data.value.map((item, index) => {

      const obj = Object.create(null);

      obj.id = index + 1;

      cols.forEach((k, i) => {
        obj[k] = !item[i] ? null : item[i];
      })

      return obj;
    })

  }
}, {
  deep: true
});


/* 关闭窗口 */
const close = () => {
  emit('update:visible', false);
}

/* 添加数据 */
const add_data = () => {

  const item = Object.create(null);

  cols.forEach(i => {
    item[i] = null;
  })

  /* 主键id */
  if (!dataSource.value || dataSource.value.length === 0) {
    item.id = 1;
  } else {
    item.id = dataSource.value[dataSource.value.length - 1].id + 1;
  }

  dataSource.value.push(item)
}


/**
 * 删除指定的数据
 * */
const del_col = (index) => {
  if (dataSource.value[index]) {
    dataSource.value.splice(index, 1);
  }
}

/**
 * 提交数据
 * */
const submit = () => {

  const account_list = []; //账号密码列表

  for (let item of dataSource.value) {

    const need_data = []; //需要的数据

    /* 拼装数据 */
    for (let i of cols) {

      if (!item[i]) {
        load.error("所有选项都必须填写！");
        return;
      }

      if (item[i].indexOf(' , ') > -1) {
        load.error("检测到特殊字符【 , 】，请删除后重试！");
        return;
      }

      need_data.push(item[i]);
    }

    account_list.push(need_data.join(' , '));
  }

  update("update:value", account_list.join("\n"));
  close(); //关闭窗口

}

</script>


<style scoped lang="scss">

.container {
  max-height: 50vh;
  box-sizing: border-box;
  width: 100%;
  padding: 0 12px 12px 15px;
  overflow: hidden;
  overflow-y: auto;
  @include scroll-bar();
}

.footer {
  @include flex-center;
  justify-content: space-between;

  .info {
    font-size: $font-size-5;
  }
}

</style>
