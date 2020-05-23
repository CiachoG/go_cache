---
typora-copy-images-to: assets
---

跟着[极客兔兔](https://geektutu.com/post/geecache-day1.html)七天实现分布式缓存，笔记和心得

**分布式缓存关注点**：

- 资源的控制
- 淘汰策略：FIFO、LFU、LRU
- 并发
- 分布式通信

**第一天**：

- 实现淘汰策略LRU算法：

  - 数据结构：

    ![implement lru algorithm with golang](assets/lru.jpg)

    map用于查，双向list用于新增和删除，lru把使用到的节点移动到链表尾，头结点就是最近最少使用的节点。代码中约定front为队尾，back为队首

  - 代码实现：

  - 



