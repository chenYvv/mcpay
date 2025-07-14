@echo off


SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

CD .\cmd\api && go build -ldflags "-s -w" -o app .

@REM CD ..\ && docker build -t gintest:dev .

CD ..\..\

docker build -t mcpay_api:1 .
@REM docker login --username=web@kaayou registry.cn-hangzhou.aliyuncs.com --password=Kaayou@0322
@REM docker push registry.cn-hangzhou.aliyuncs.com/kaayou_cy_test/manage_server:%ver%


docker login uhub.service.ucloud.cn -u liuyuhang@whu.edu.cn -p xugang83647762

docker tag mcpay_api:1 uhub.service.ucloud.cn/bingbinggod/testImg

docker push uhub.service.ucloud.cn/bingbinggod/testImg

pause

