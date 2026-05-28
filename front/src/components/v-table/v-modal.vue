<template>
  <a-modal
      :centered="centered"
      :class="modelClass"
      :maskClosable="false"
      :title="title"
      :open="visible"
      :width="width"
      class="antd-reset-batch-modal"
      destroy-on-close
      @cancel="close"
      @ok="ok">

    <template #footer>

      <template v-if="!footer"></template>
      <template v-else>
        <div class="footer">
          <span class="info">
              <slot name="info"></slot>
          </span>
          <a-space>
            <v-button @click="close">{{ cancelText }}</v-button>
            <v-button :loading="loading" type="primary" @click="ok">{{ okText }}</v-button>
          </a-space>
        </div>
      </template>

    </template>

    <div ref="table" :class="add_class" :style="{maxHeight:maxHeight}" class="table container"
         @scroll="e=>{e.stopPropagation()}">

      <!-- 表单 -->
      <slot name="batch"/>

      <a-table
          v-model:expandedRowKeys="expandedRowKeys"
          :columns="column"
          :dataSource="data.data.data"
          :loading="loadingTable"
          :pagination="!pagination?pagination:paginate"
          :rowKey="key"
          :rowSelection="rowSelection===false?null:(rowSelection === true?selectConfig:rowSelection)"
          :scroll="scroll"
          :size="size"
          :sticky="sticky"
          @change="loadData"
          @resizeColumn="handleResizeColumn">

        <!--客户展开的详情-->
        <template v-if="$slots.expandedRowRender" #expandedRowRender="{ record }">
          <slot :record="record" name="expandedRowRender"/>
        </template>

        <template #headerCell="{ column }">
          <slot :column="column" name="headerCell"/>
        </template>

        <template #bodyCell="{text,record, index, column}">
          <template v-if="(text === null || text === '') && !column.empty ">
            <span>-</span>
          </template>
          <slot v-else :column="column" :index="index" :record="record" :text="text" name="bodyCell"/>
        </template>

        <template #summary>
          <slot v-if="data.data.data.length > 0" :data="data.data" name="summary"/>
        </template>

      </a-table>
    </div>
  </a-modal>
</template>

<script setup>

import initialize from './v-table'
import mask from "@/common/mask";

const props = defineProps({
  tid: {
    required: false,
    default: 'table'
  },
  okText: {
    required: false,
    type: String,
    default: "确定"
  },
  cancelText: {
    required: false,
    type: String,
    default: "取消"
  },
  modelClass: {
    required: false
  },
  loading: {
    required: false,
    default: false
  },
  width: {
    required: false,
    default: "520px"
  },
  title: {
    required: true
  },
  add_class: {
    required: false
  },
  maxHeight: {
    default: "65vh"
  },
  footer: {
    default: true,
    required: false
  },
  centered: {
    required: false,
    default: false
  },
  init: {
    required: false
  }, //接口
  visible: {
    required: true
  }, //显示
  dataSource: {
    required: false,
    default: false
  },
  pagination: {
    required: false,
    default: {}
  },
  sticky: {
    required: false,
    default: {offsetHeader: 60}
  },
  rowsKey: {
    required: false,
    default: "id"
  },
  scroll: {
    required: false,
    default: () => ({x: '100%'})
  },
  size: {
    required: false,
    default: "large"
  },
  rowSelection: {
    required: false,
    default: false
  }
});

/* 防止滚动 */
mask(props, 2);

/* 自定义事件 */
const emit = defineEmits(['update:visible', 'ok']); //自定义事件

/* 关闭弹窗 */
const close = () => {
  emit('update:visible', false)
};

/* 点击确定 */
const ok = () => {
  emit('ok')
};


let {
  column,
  data,
  key,
  loadData,
  loadingTable,
  paginate,
  handleResizeColumn,
  selectConfig,
  expandedRowKeys
} = initialize(props);


</script>


<style lang="scss" scoped>

.container {
  box-sizing: border-box;
  width: 100%;
  padding: 20px 20px 0 25px;
  overflow: hidden;
  overflow-y: auto;
  min-height: 60vh;
  @include scroll-bar();
}

.footer {
  @include flex-center;
  justify-content: space-between;

  .info {
    font-size: $font-size-5;
  }
}

.table {

  :deep(.ant-table) {

    font-family: DINPro, PingFang SC, -apple-system, BlinkMacSystemFont, Segoe UI, Hiragino Sans GB, Microsoft YaHei, Helvetica Neue, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", Segoe UI Symbol !important;

    .ant-table-container {

      .ant-table-content, .ant-table-body {
        @include scroll-bar;

        /* 滚动条样式 */
        &::-webkit-scrollbar {
          width: 10px !important;
          height: 10px !important;
        }
      }

      .ant-table-thead {


        & > tr {

          & > th {
            position: relative;
            color: rgba(0, 0, 0, .85);
            font-weight: 500;
            text-align: left;
            background: $table-td;
            border-bottom: 1px solid #f0f0f0;
            -webkit-transition: background .3s ease;
            transition: background .3s ease;
            @include overflow;
          }

          & > .ant-table-selection-column {
            text-align: center;
          }
        }
      }

      .table-striped {
        background-color: #f6f9f9;
      }

      td, th {
        font-size: $font-size-3 !important;
      }


    }

    .ant-table-column-sort {
      background: unset;
    }

  }
}

</style>
