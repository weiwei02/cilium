# 为集群更新容器网络 v2 组件的版本
# 使用示例 tools/operating/helm-upgrade.sh output/z85m1xdk/node-ready-values-z85m1xdk.yaml
#!/bin/bash

# 检查当前工作路径是否以 cilium 结尾
current_dir=$(basename "$PWD")
if [[ $current_dir != "cilium" ]]; then
  echo "错误：当前工作路径不是以 cilium 结尾。"
  exit 1
fi

if [ ! -f "$1" ]; then
  echo "错误：文件 '$1' 不存在。"
  exit 1
fi

resources=(
    # crd 资源

    # 控制器资源
    "Deployment/cilium-operator"
    "DaemonSet/cilium"
    # "Service/cilium"

    # 配置文件
    "ConfigMap/cilium-config"
    # "ConfigMap/cni-config-template"

    # 权限
    "ServiceAccount/cilium"
    "ServiceAccount/cilium-operator"
    "ClusterRole/cilium"
    "ClusterRole/cilium-operator"
    "ClusterRoleBinding/cilium"
    "ClusterRoleBinding/cilium-operator"
)

# 为每个资源执行增加 label 和 annotation 的操作
for resource in "${resources[@]}"; do
  kubectl -n kube-system label $resource app.kubernetes.io/managed-by=Helm --overwrite
  kubectl -n kube-system annotate $resource meta.helm.sh/release-namespace=kube-system meta.helm.sh/release-name=cilium --overwrite
  echo "update $resource success"
done


kubectl -n kube-system delete daemonset cilium 
kubectl -n kube-system delete deployment cilium-operator

# 执行安装脚本
# 例如 tools/operating/helm-upgrade.sh output/z85m1xdk/node-ready-values-z85m1xdk.yaml
helm upgrade -n kube-system --install cilium baidu/cilium-charts -f $1 