docker build -t browsify-js .
docker run \
  --volume browsify-js:/data \
  --rm \
  browsify-js \
  bash -c 'cp -r /web/build/* /data/.'
