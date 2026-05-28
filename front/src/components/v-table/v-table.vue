<template>
  <div ref="table" class="table">
    <!-- :row-class-name="(_record, index) => (index % 2 === 1 ? 'table-striped' : null)" -->

    <a-table
      :columns="column"
      :data-source="data.data.data"
      :loading="loadingTable"
      :pagination="!pagination?pagination:paginate"
      :row-key="key"
      :row-selection="rowSelection === false ? null : (rowSelection === true ? selectConfig : rowSelection)"
      :scroll="scroll"
      :sticky="sticky"
      :size="size"
      table-layout="fixed"
      v-model:expanded-row-keys="expandedRowKeys"
      @change="loadData"
      @resize-column="handleResizeColumn"
    >
      <!--客户展开的详情-->
      <template v-if="$slots.expandedRowRender" #expandedRowRender="{ record }">
        <slot :record="record" name="expandedRowRender" />
      </template>

      <template #headerCell="{ column }">
        <slot :column="column" name="headerCell" />
      </template>

      <template #bodyCell="{text,record, index, column}">
        <template v-if="(text === null || text === '' || (column.dataIndex && !(column.dataIndex in record))) && !column.empty ">
          <span>-</span>
        </template>
        <slot v-else :column="column" :index="index" :record="record" :text="text" name="bodyCell" />
      </template>

      <template #summary>
        <slot v-if="data.data.data.length > 0" :data="data.data" name="summary" />
      </template>
    </a-table>
  </div>
</template>

<script setup>
import initialize from './v-table'

/**
 * 表格可选修改
 * 分页取消显示
 * 插槽增加text
 * 表单的可选可通过keys控制
 * 新增自定义key值
 * 定义表头和初始化数据的接口
 * */
const props = defineProps({
  tid: {
    required: false,
    default: 'table'
  },
  init: {
    required: false
  },
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
  /**
   * false 不显示可选
   * true 显示内置的可选
   * 对象 ,显示自定义的可选
   * */
  rowSelection: {
    required: false,
    default: false
  }
});


let {
  column,
  data,
  key,
  /*loaded,*/
  loadData,
  loadingTable,
  paginate,
  handleResizeColumn,
  selectConfig,
  expandedRowKeys,
  /*layout*/
} = initialize(props);


</script>


<style lang="scss" scoped>
.table {

  /*:deep(.ant-table-header) {
    table {
      table-layout: v-bind(layout) !important;
    }
  }*/

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
            text-overflow: unset;
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