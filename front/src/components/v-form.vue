<template>
  <a-modal
    :centered="centered"
    :class="{'antd-mobile-modal':$setting.isMobile}"
    :mask-closable="false"
    :title="title"
    :open="visible"
    :width="width"
    class="antd-reset-modal"
    destroy-on-close
    @cancel="close"
    @ok="modify"
  >
    <template #footer>
      <template v-if="!footer" />
      <template v-else>
        <div class="footer">
          <span class="info">
            <slot name="info" />
          </span>
          <a-space>
            <v-button @click="close">{{ cancelText }}</v-button>
            <v-button :loading="loading" type="primary" @click="modify">{{ okText }}</v-button>
          </a-space>
        </div>
      </template>
    </template>

    <div :class="add_class" :style="{maxHeight:maxHeight}" class="container" @scroll="e=>{e.stopPropagation()}">
      <v-batch
        ref="formRef"
        :check="check"
        :data="dataSource.data"
        :form="form"
        :grid="grid"
        :gutter="gutter"
        :label-span="labelSpan"
        :layout="layout"
        :name="name"
        :offset="false"
        :show-border="showBorder"
        :union="union"
        label="fixed"
      />

      <slot name="extra" />
    </div>
  </a-modal>
</template>

<script setup>

import load from "@/common/load";
import axios from "axios";
import {inject, provide, ref} from "vue";
import {switchForm} from "@/common";
import {cloneDeep} from "lodash";
import qs from "qs";

/* 属性 */
const props = defineProps(
    {
      tid: {
        required: false,
        default: "table"
      },
      name: {
        required: false,
        type: String,
        default: "form"
      },
      okText: {
        required: false,
        type: String,
        default: "确定"
      },
      gutter: {
        type: [Number, Object, Array],
        default: 16
      },
      cancelText: {
        required: false,
        type: String,
        default: "取消"
      },
      title: {
        required: true
      }, //居中
      add_class: {
        required: false
      }, //样式
      centered: {
        required: false,
        default: false
      },
      layout: {
        required: false,
        default: "horizontal"
      },
      grid: {
        required: false,
        default: {xxl: 24, xl: 24, lg: 24, md: 24, sm: 24, xs: 24}
      },
      footer: {
        default: true,
        required: false
      },
      message: {
        required: false
      },
      init: {
        required: false
      },
      visible: {
        required: true
      },
      form: {
        required: true
      },
      dataSource: {
        required: true
      },
      showBorder: {
        default: false
      },
      union: {
        default: true
      },
      /* 提交后的回调 */
      after: {
        required: false,
        default: null
      },
      /* 提交前的回调 */
      before: {
        required: false,
        default: null
      },
      check: {
        required: false,
        default: null
      },
      /* 自定义表单提交 */
      submit: {
        required: false,
        default: null
      },
      width: {
        required: false,
        default: 500
      },
      labelSpan: {
        required: false,
        default: false
      },
      hasTable: {
        required: false,
        default: true
      },
      loading: {
        required: false,
        default: false
      },
      dialog: {
        required: false,
        default: 'success'
      },
      format: {
        default: "YYYY-MM-DD"
      },
      maxHeight: {
        default: "65vh"
      },
      method: {
        default: "POST"
      },
      header: {
        default: null
      },
      filter: {
        default: null
      }
    });

/* 自定义事件 */
const emit = defineEmits(['update:visible']); //自定义事件

/* 关闭弹窗 */
const close = () => {
  emit('update:visible', false)
};

/* 注入表格对象 */
const table = inject(props.tid, {
  loadData: () => {
  }
});

/* 表单dom */
const formRef = ref(null);

/**
 * 表单
 * */
const modify = async () => {

  /* 待提交的数据 */
  let data = cloneDeep(props.dataSource.data);

  /* 触发前的回调 */
  if (props.before) {
    /* 判断是否继续 */
    if (!await props.before(data)) {
      return false;
    }
  }


  /* 没有接口直接终止 */
  if (!props.submit && !props.init) return;

  /**
   * 表单验证
   * */
  formRef.value.$refs.formRef
      .validate()
      .then(async () => {

        /**
         * 自定义上传
         * */
        if (props.submit) {

          if (await props.submit(switchForm(data, props.format))) {

            close(); //关闭弹窗

            /* 有表格需要刷新 */
            if (props.hasTable) {
              table.loadData(); //重新加载表格数据
            }

          }
          return;
        }

        /* 显示加载效果 */
        load.loading("请求中...");

        /* 请求对象 */
        const config = {
          method: props.method,
          url: props.init,
          data: switchForm(data, props.format)
        }

        /* 请求链接 */
        if (props.method === 'get') {
          config.url = config.url + "?" + qs.stringify(config.data);
          delete config.data; //删除POST参数
        }


        /* 如果自定义 请求头 */
        if (props.header) {
          config.headers = props.headers;
        }

        /* 数据过滤 */
        if (props.filter) {
          config.data = props.filter(config.data);
        }

        /* 开始请求 */
        let requestPromise;
        if (typeof props.init === 'function') {
            requestPromise = Promise.resolve(props.init(config.data || switchForm(data, props.format))).then(res => {
                let parsed = res;
                if (typeof res === 'string') {
                    try { parsed = JSON.parse(res); } catch (e) {}
                }
                return {
                    data: {
                        code: parsed && parsed.code !== undefined ? parsed.code : 1,
                        errMsg: parsed && parsed.errMsg ? parsed.errMsg : '',
                        data: parsed && parsed.data ? parsed.data : parsed
                    }
                };
            }).catch(e => {
                return { data: { code: 0, errMsg: e.message || String(e) } };
            });
        } else {
            requestPromise = axios.request(config);
        }

        requestPromise
            .then((res) => {
              /**
               * 判断请求结果
               * */
              if (res.data.code) {

                /* 是否自定义提交提示 */
                if (props.message) {
                  load[props.dialog](props.message + '！');
                } else {
                  load[props.dialog](res.data.errMsg);
                }

                /* 关闭弹窗 */
                close();

                /* 有表格需要刷新 */
                if (props.hasTable) {
                  table.loadData(); //重新加载表格数据
                }

                /* 提交后的回调 */
                if (props.after) {
                  props.after(res.data); //触发钩子
                }


              } else {
                /* 弹出错误原因 */
                load.error(res.data.errMsg);
              }
            }).catch((e) => {
          /* 弹出错误原因 */
          load.error(e.message);
        }).finally(() => {
          /* 关闭加载效果 */
          load.loaded();
        });

      }).catch((e) => {
    console.log(e);
    load.error("请按照要求填写数据！")
  })

}

/* 注入数据 */
provide('dataSource', props.dataSource);

</script>

<style lang="scss" scoped>

.container {
  box-sizing: border-box;
  width: 100%;
  padding: 20px 20px 20px 25px;
  overflow: hidden;
  overflow-y: auto;
  @include scroll-bar();

  :deep(.horizontal) {
    .ant-form-item {
      margin-bottom: 1.6rem;
    }
  }
}

.footer {
  @include flex-center;
  justify-content: space-between;

  .info {
    font-size: $font-size-5;
  }
}

</style>
