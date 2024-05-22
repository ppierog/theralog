UID=$(id -u)
GID=$(id -g)
IMAGE_NAME="alpine-go"
HOST_NAME="${IMAGE_NAME}-container"
docker run --rm -it \
    -p 8081:8080 \
    --user $UID:$GID \
    --env UID=$(id -u) \
    --env GID=$(id -g) \
    --workdir="/home/$USER" \
    --volume="/etc/group:/etc/group:ro" \
    --volume="/etc/passwd:/etc/passwd:ro" \
    --volume=/$HOME:$HOME \
    --hostname="$HOST_NAME" \
    --name="$HOST_NAME" \
"$IMAGE_NAME"