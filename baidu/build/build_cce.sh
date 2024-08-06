# export ARCH=amd64 

# 开启多架构镜像构建
export DOCKER_BUILDX=1

export DOCKER_REGISTRY=registry.baidubce.com/cce-plugin-dev
export DOCKER_DEV_ACCOUNT=cilium
export DOCKER_IMAGE_TAG=$(cat VERSION)

# 构建 agent 镜像
make dev-docker-image
dev_image=${DOCKER_REGISTRY}/cilium/cilium-dev:${DOCKER_IMAGE_TAG}
cilium_image=${DOCKER_REGISTRY}/cilium/cilium:${DOCKER_IMAGE_TAG}
docker tag ${dev_image} ${cilium_image}
docker push ${cilium_image}
# 构建 operator 镜像
# make docker-operator-generic-image

# 生产镜像
# make docker-images-all