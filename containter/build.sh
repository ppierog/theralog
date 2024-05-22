GO_LINK="https://go.dev/dl/go1.21.6.linux-amd64.tar.gz"
GO_TAR=$(echo "$GO_LINK" | rev | cut -d "/" -f1 | rev)
IMAGE_NAME="alpine-go"

docker build --build-arg GO_LINK="$GO_LINK" --build-arg GO_TAR="$GO_TAR" ./ -t $IMAGE_NAME