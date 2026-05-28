<template>
  <a-radio-group v-model:value="selected" button-style="solid">
    <template v-for="(item,index) of options" :key="index">
      <a-radio-button :value="item.value">{{ item.label }}</a-radio-button>
    </template>
  </a-radio-group>
</template>

<script setup>
import {computed} from "vue";

const props = defineProps(['value', 'options']);
const emit = defineEmits(['update:value', "change"]);

const selected = computed({
  get() {
    return props.value;
  }, set(value) {
    emit('update:value', value);
    emit("change", value);
  }
})


</script>

<style lang="scss" scoped>

</style>