docker build -t flux-panel:latest .
#docker pull registry.cn-hangzhou.aliyuncs.com/nqc/arkoselabs_token_api.v2:latest
docker tag flux-panel:latest 24802117/flux-panel:latest
docker push 24802117/flux-panel:latest

