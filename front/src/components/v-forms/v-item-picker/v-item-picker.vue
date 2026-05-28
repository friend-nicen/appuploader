<template>

  <a-textarea
      :value="value"
      :readOnly="true"
      :placeholder="placeholder"
      :autoSize="autoSize"
      :required="required"
      @click="visible=true"
      :allowClear="false"/>

  <a-form-item-rest>
    <v-item-modal :unique="unique" :columns="columns" v-model:visible="visible"/>
  </a-form-item-rest>

</template>


<script setup>
import {provide, ref, watch} from "vue";
import VItemModal from './v-item-modal.vue'


const props = defineProps({
  /* 只读 */
  readOnly: {
    required: false,
    default: true
  },
  /* 文本域大小 */
  autoSize: {
    required: false
  },
  /* 必选 */
  required: {
    required: false
  },
  /* 可清除 */
  allowClear: {
    required: false
  },
  /* 占位符 */
  placeholder: {
    required: false
  },
  /* 值 */
  value: {
    required: true
  },
  /* 表格标识符 */
  unique: {
    required: false,
    default: "Account-List"
  },
  /* 表格默认列 */
  columns: {
    required: false,
    default: () => {
      return [
        {
          dataIndex: "id",
          title: "序号",
          width: 50
        },
        {
          dataIndex: "user",
          title: "用户名",
          empty: true,
          width: 150
        },
        {
          dataIndex: "password",
          title: "密码",
          empty: true,
          width: 150
        },
        {
          title: "操作",
          fixed: "right",
          width: 60
        }
      ]
    }
  },
  /* 数据预处理 */
  filter: {
    required: false,
    default: () => {
      return (data) => {
        const items = !data.value ? [] : data.value.split("\n");
        return items.map((item) => {
          return item.split(" , ");
        })
      }
    }
  }
});

const emit = defineEmits(['update:value']);
const visible = ref(false); //是否显示表格录入框

/* 数据集 */
const dataSource = ref(props.filter(props));
provide("update", emit); //更新值
provide("item-picker-data", dataSource); //更新值


/* 检测并改变 */
watch(() => props.value, () => {
  dataSource.value = props.filter(props);
}, {
  immediate: true
})


</script>


<style scoped lang="scss">

</style>