import {getBase64} from "@/common";
import {onUnmounted, ref} from 'vue'
import load from "@/common/load";
import {Upload} from "ant-design-vue";

export default function (props) {


    const previewVisible = ref(false); //图片预览的模态框
    const previewImage = ref(null);  //图片预览的图
    const fileList = ref([]); //文件列表

    /**
     * 图片预览
     * */
    const handlePreview = async (file) => {
        if (!file.url && !file.preview) {
            file.preview = await getBase64(file.originFileObj);
        }
        previewImage.value = file.url || file.preview;
        previewVisible.value = true;
    };


    /**
     * 附件删除
     * */
    const fileRemove = async () => {
        /*判断是否是已经上传的文件*/
        return true;
    }


    /**
     * 文件检测
     * */
    const beforeUpload = (file) => {

        if (!props.size) return false;


        /**
         * 判断文件大小是否超出限制
         * */
        if ((file.size / 1024) > props.size) {
            load.error("单个文件大小最大不能超过" + Math.round(props.size / 1024, 2) + "M!");
            return Upload.LIST_IGNORE;
        }


        return false;

    }


    /**
     * 粘贴图片上传
     * */
    const onPaste = (e) => {

        /* 判断是否有数据 */
        if (e.clipboardData && e.clipboardData.items && e.clipboardData.items.length) {

            const clipboardData = e.clipboardData;//谷歌

            let success = 0;//成功粘贴的数据

            /* 遍历文件对象 */
            for (let item of clipboardData.items) {
                /* 是不是图片 */
                if (item.type.indexOf("image") !== -1) {


                    const file = item.getAsFile(); //读取文件

                    if (!file) {
                        continue;
                    }

                    if (fileList.value.length < props.number) {
                        fileList.value.push({
                            originFileObj: file,
                            thumbUrl: URL.createObjectURL(file)
                        });
                        success++; //计数
                    }
                }
            }

            /* 结果判断 */
            if (success > 0) {
                load.success("图片粘贴成功！");
            }

        }
    }

    /* 粘贴图片 */
    window.addEventListener('paste', onPaste); //粘贴事件
    onUnmounted(() => {
        window.removeEventListener('paste', onPaste); //粘贴事件
    })

    return {
        previewVisible,
        previewImage,
        handlePreview,
        fileList,
        fileRemove,
        beforeUpload
    }
}
