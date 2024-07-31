# 1.15
## 1.15.6-baidu-2407294
1. [Feature] 新增网络编程延迟度量，支持 cilium 主动暴漏网络编程延迟
2. [Feature] 支持 datapath v2 Pod 网络加速
3. [Feature] 默认开启 cilium host firewall 网络策略防火墙
4. [Feature] 支持 cilium nodeport 端口监听保护
2. [Feature] 增加 CCE 镜像构建和推送
3. [Optimize] 移除 cilium 节点污点管理功能，避免network-operator 和 cilium-operator 竞争修改 status 时引起的 apiserver 并发循环修改
4. [Optimize] 默认关闭 hubble relay，避免 hubble relay 占用过多资源