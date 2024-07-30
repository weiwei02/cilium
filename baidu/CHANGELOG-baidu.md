# 1.15
## 1.15.6-baidu-2407294
1. [Feature] 新增网络编程延迟度量，支持 cilium 主动暴漏网络编程延迟
2. [Feature] 增加 CCE 镜像构建和推送
3. [Optimize] 移除 cilium 节点污点管理功能，避免network-operator 和 cilium-operator 竞争修改 status 时引起的 apiserver 并发循环修改