<template>
  <div class="upload-up">

    <a-upload
        v-model:file-list="fileList"
        :accept="accept"
        :before-upload="beforeUpload"
        list-type="picture-card"
        multiple
        :max-count="number"
        @preview="handlePreview"
        @remove="fileRemove">

      <div class="flex" v-if="fileList.length < number">
        <v-icon icon="PlusOutlined"/>
        <span class="info">{{ info }}</span>
      </div>

    </a-upload>

    <a-modal :footer="null" title="图片预览" :open="previewVisible" @cancel="previewVisible=false">
      <div class="image-container">
        <a-image class="imamge" :src="previewImage" alt="example"/>
      </div>
    </a-modal>
  </div>
</template>

<script setup>

import init from "./v-upload";
import {watch} from "vue";


const props = defineProps({
  info: {
    required: false,
    default: "选择文件"
  },
  size: {
    default: 10240,
    required: false
  },
  files: {
    required: true
  },
  number: {
    default: 8,
    required: false
  },
  accept: {
    required: false,
    default: '.png,.jpg,.jpeg,.gif'
  }
});

const emits = defineEmits(['update:files']);

let {
  previewVisible,
  previewImage,
  handlePreview,
  fileList,
  fileRemove,
  beforeUpload
} = init(props);


/* 监视文件的改变 */
watch(fileList, () => {

  const files = []; //源文件

  for (let i of fileList.value) {
    files.push(i.originFileObj);
  }

  emits("update:files", files);//更新已选文件

}, {
  deep: true
})

</script>

<style lang="scss" scoped>


.upload-up {

  min-height: 112px;

  .flex {
    @include flex-center;
    justify-content: center;
    flex-direction: column;
  }

  .info {
    font-size: $font-size-6;
    color: $text-color-4;
    margin-top: 5px;
  }
}


.image-container {
  @include flex-center;
  justify-content: center;
  width: 100%;

  .image {
    max-height: 50vh;
    width: 100%
  }
}
</style>
