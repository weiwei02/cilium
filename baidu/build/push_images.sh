#! /bin/bash
TAG=1.15.6-baidu-240726
IMAGES=("registry.baidubce.com/cce-plugin-dev/cilium/clustermesh-apiserver" 
    "registry.baidubce.com/cce-plugin-dev/cilium/hubble-relay" 
    "registry.baidubce.com/cce-plugin-dev/cilium/cilium-dev" 
    "registry.baidubce.com/cce-plugin-dev/cilium/docker-plugin" 
    "registry.baidubce.com/cce-plugin-dev/cilium/cilium" 
    "registry.baidubce.com/cce-plugin-dev/cilium/operator" 
    "registry.baidubce.com/cce-plugin-dev/cilium/operator-generic" 
)

for img in "${IMAGES[@]}"; do
    docker push $img:$TAG
done

OTHER_IMAGES=("quay.io/cilium/certgen:V0.1.13" 
"quay.io/cilium/hubble-ui-backend:v0.13.1"
"quay.io/cilium/hubble-ui:v0.13.1"
"quay.io/cilium/cilium-envoy:v1.29.7-39a2a56bbd5b3a591f69dbca51d3e30ef97e0e51"
"quay.io/cilium/cilium-etcd-operator:v2.0.7"
"quay.io/cilium/startup-script:c54c7edeab7fde4da68e59acd319ab24af242c3f"
"ghcr.io/spiffe/spire-server:1.8.5"
"ghcr.io/spiffe/spire-agent:1.8.5"
"docker.io/library/busybox:1.36.1"
)

for quay_image in "${OTHER_IMAGES[@]}"; do
    cce_img=$(echo $quay_image | sed 's/ghcr.io\/spiffe/registry.baidubce.com\/cce-plugin-dev\/cilium/g')
    cce_img=$(echo $quay_image | sed 's/docker.io\/library/registry.baidubce.com\/cce-plugin-dev\/cilium/g')
    cce_img=$(echo $quay_image | sed 's/quay.io/registry.baidubce.com\/cce-plugin-dev/g')
    cce_pro=$(echo $quay_image | sed 's/quay.io/registry.baidubce.com\/cce-plugin-pro/g')
    docker pull $quay_image
    docker tag $quay_image $cce_img
    docker push $cce_img
done