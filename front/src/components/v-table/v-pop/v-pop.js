import load from "@/common/load";
import {computed, inject, toRef} from "vue";
import {cloneDeep} from "lodash";


export default function (props) {


    /* 表格对象 */
    const table = inject(props.tid);
    const usedColumns = toRef(table, "usedColumns");
    const defaultColumns = toRef(table, "defaultColumns");
    const sort = toRef(table, "sort"); //排序

    /**
     * 如果有一个不显示
     * 代表没有全选
     * */
    const allChecked = computed({
        get() {
            for (let i of usedColumns.value) {
                if (!i.display) {
                    return false;
                }
            }
            return true;
        },
        set(newValue) {
            if (newValue) {
                for (let i of usedColumns.value) {
                    i.display = true;
                }
            } else {
                for (let i of usedColumns.value) {
                    if (!i.default) {
                        i.display = false;
                    }
                }
            }
        }
    })


    /**
     * 重置列
     * */
    const reset = () => {
        usedColumns.value = cloneDeep(defaultColumns.value);
        load.success("重置成功！")
    }


    return {
        allChecked,
        usedColumns,
        reset,
        sort,
        table
    }

}