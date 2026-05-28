<template>
  <a-modal
    :open="visible"
    :mask-closable="false"
    title="图标选择"
    :width="width"
    class="antd-reset-modal"
    @cancel="close"
  >
    <template #footer />

    <div class="container" v-bind="containerProps">
      <ul class="icon-list" v-bind="wrapperProps">
        <template v-for="(icons,index) of list" :key="index">
          <template v-for="item of icons.data" :key="item">
            <li class="icon" :class="{active:item === icon}" @click="select(item)">
              <v-icon :icon="item" size="30px" />
            </li>
          </template>
        </template>
      </ul>
    </div>
  </a-modal>

  <div class="icon-picker">
    <v-icon :icon="icon" size="25px" />
    <a-input :required="required" :placeholder="placeholder" :allow-clear="allowClear" v-model:value="icon" read-only />
    <v-button type="primary" @click="visible=true">选择</v-button>
  </div>
</template>

<script setup>

import {ref, watch} from "vue";
import * as $icon from '@ant-design/icons-vue';
import {useVirtualList} from '@vueuse/core'

/* 排除指定的图标 */
const ignore = ['create', 'getTwoToneColor', 'setTwoToneColor', 'default', 'createFromIconfontCN'];

/* 内置图标库 */
const icons = ref([]);

let temp = [];

/* 针对虚拟滚动进行重组 */
for (let i in $icon) {
  if (typeof i === 'string' && ignore.indexOf(i) === -1) {
    temp.push(i);
  }
  if (temp.length === 10) {
    icons.value.push(temp);
    temp = []; //置空
  }
}

/* 创建虚拟滚动的容器 */
const {list, containerProps, wrapperProps} = useVirtualList(
    icons,
    {
      itemHeight: 60,
    },
)


/* 属性 */
const props = defineProps(
    {
      value: {
        required: true
      },
      width: {
        required: false,
        default: 800
      },
      required: {
        required: false,
        default: true
      },
      placeholder: {
        required: false,
        default: "请选择图标代码"
      },
      allowClear: {
        required: false,
        default: true
      }
    });

/* 被选icon */
const icon = ref(props.value);
const visible = ref(false);


/* 监视属性 */
watch(() => props.value, () => {
      icon.value = props.value;
    },
    {
      deep: true
    });

//自定义事件
const emit = defineEmits(['update:value']); //自定义事件

//关闭弹窗
const close = () => {
  visible.value = false;
};

/**
 * 确定提交
 * */
const select = (icon) => {
  emit("update:value", icon)
  close(); //关闭窗口
}

</script>


<style scoped lang="scss">

.container {

  max-height: 65vh;
  box-sizing: border-box;
  width: 100%;
  padding: 20px;
  overflow: hidden;
  overflow-y: auto;
  @include scroll-bar();

  .icon-list {
    display: block;
    list-style: none;

    padding: 0;
    margin: 0;

    &, * {
      box-sizing: border-box;
    }

    .icon {
      width: 10%;
      padding: 15px;
      cursor: pointer;
      display: inline-flex;
      align-items: center;
      flex-direction: column;
      border-radius: 12px;

      span {
        font-size: $font-size-3;
        color: $text-color-3;
      }

      &:hover {
        background-color: #f5f6f7;
      }
    }

    .active {
      background-color: rgba($primary-color, 0.45);
    }
  }
}

.icon-picker {
  display: flex;
  gap: 15px;
  width: 100%;
  align-items: center;
}

</style>
