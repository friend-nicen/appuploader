<template>
  <span v-if="loaded">{{ prefix }}{{ value }}{{ suffix }}</span>
  <loading-outlined v-else/>
</template>

<script setup>
import {LoadingOutlined} from "@ant-design/icons-vue";

/* 接口 */
import {onMounted, ref} from "vue";
import axios from "axios";
import {delay} from "lodash";

const props = defineProps({
  api: {
    required: true
  },
  data: {
    required: false,
    default: {}
  },
  prefix: {
    required: false
  },
  suffix: {
    required: false
  }
});

/* 数据 */
const value = ref(null);
const loaded = ref(false);

/* 初始化 */
onMounted(() => {
  try {
    /* 开始请求 */
    axios.post(props.api, props.data)
        .then((res) => {
          /* 判断请求结果 */
          if (res.data.code) {
            delay(() => loaded.value = true, 300);
            value.value = res.data.data;
          }
        }).catch((e) => {
      /* 弹出错误原因 */
      console.log(e)
    });
  } catch (e) {
    console.log(e)
  }
});

</script>

<style lang="scss" scoped>

</style>
