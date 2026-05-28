/**
 * @author 友人a丶
 * @date
 * 树形结构函数
 * */


/**
 * 获取树形数据结构
 * */
export function getTree(data, parent = 1, key = "key", callback = null) {

    if (!Array.isArray(data)) return [];
    if (data.length === 0) return [];

    const Tree = []; //树

    /**
     * 遍历所有父级id为指定id的元素
     * */
    for (let i of data) {

        if (i.pid === parent) {

            /* 获取当前遍历的数据的所有子集 */
            const children = getTree(data, i[key], key, callback);

            if (children.length > 0) {
                i.children = children;
            }

            if (callback) {
                i = callback(i);
            }

            Tree.push(i);

        }
    }

    return Tree;

}


/**
 * 获取所有带子集的节点
 * @param data
 * @param key
 */
export function getNode(data, key = "key") {

    if (!Array.isArray(data)) return [];
    if (data.length === 0) return [];

    let Node = [];

    for (let count = 0; count < data.length; count++) {

        let i = data[count]; //数据对象

        /* 跳过非叶节点 */
        if (!i.children) continue;

        Node.push(i[key]);//新增节点
        let nodes = getNode(i.children, key);

        if (nodes && nodes.length > 0) {
            Node = Node.concat(nodes);
        }
    }

    return Node;
}