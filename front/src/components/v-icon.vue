<script>
import {h} from 'vue'
import * as $icon from '@ant-design/icons-vue';

/* 自动引入antd icon图标 */
export default {
  props: {
    icon: {
      default: "BulbOutlined"
    },
    size: {
      default: '15px'
    }
  },
  setup(props) {
    /* 判断是icon还是antd的图标 */
    if (props.icon && props.icon.indexOf('icon-') !== -1) {
      return () => h(
          'i',
          {
            class: "iconfont " + props.icon,
            style: {fontSize: props.size}
          });
    } else {
      return () => h(
          !$icon[props.icon] ? $icon['BulbOutlined'] : $icon[props.icon],
          {
            style: {fontSize: props.size}
          });
    }
  }
}

</script>
