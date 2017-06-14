# What
提供ACM-ICPC赛时服务。服务项：
1. 榜单服务
2. 气球服务
3. 打印机服务
4. 后台管理服务

# Usage

下面是大概的使用方法，具体到代码请参照后面的 Usage Example
1. 安装git, golang, sqlite3
2. 下载代码
```bash
cd ${WORK_DIR} && git clone git@github.com:shjwudp/ACM-ICPC-api-service.git
```
3. 获取go代码依赖  
```bash
cd ${ACM_ICPC_api_service} && go get ./...
```
4. 运行
```bash
cd ${ACM_ICPC_api_service} && go run main.go --config conf.json
```
conf.json说明：
```json
{
    "Server": {
        "JWTSecret": "terceSTWJ",       \\密码，随机一个64位的随机密码串就好
        "Port": ":8080",                \\端口号，不推荐改变
        "Admin": {
            "Account": "admin",         \\管理员账号
            "Password": "ElPsyCongroo"  \\管理员密码
        },
        "IsTestMode": true              \\是否为测试模式，比赛时一定要改成false，重启服务
    },
    "Storage": {
        "Dirver": "sqlite3",            \\数据库名，不得改变
        "Config": "sqlite3.db"          \\sqlite文件名
    },
    "Printer": {
        "QueueSize": 1000,              \\打印机队列长度
        "PrinterNameList": []           \\打印机名，服务器配好打印机后，将打印机名填写到这里即可
    },
    "ResultsXMLPath": "./results.xml",  \\pc2 产出的 results.xml 路径，需要正确填写
    "ContestInfo": {
        "StartTime": "2017-06-10T10:20:00+08:00",   \\比赛开始时间，格式为RFC3339
        "GoldMedalNum": 3,              \\金牌数
        "SilverMedalNum": 5,            \\银牌数
        "BronzeMedalNum": 10,           \\铜牌数
        "Duration": 18000               \\比赛持续时间（单位：秒）
    }
}
```


# Usage Example
## Fedora
保证以sudo前缀或者是管理员身份运行下列命令
1. 安装git, golang, sqlite3
```bash
dnf install git
dnf install golang
dnf install sqlite3
# 配置gopath，在~/.bashrc中写入GOPATH变量，可参照下面方法
echo "GOPATH=\"/usr/share/gocode\"" >> ~/.bashrc
source ~/.bashrc
```
2. 下载代码
```bash
WORK_DIR=/code/top/directory
cd ${WORK_DIR} && git clone git@github.com:shjwudp/ACM-ICPC-api-service.git
```
3. 获取go代码依赖
```bash
ACM_ICPC_api_service=${WORK_DIR}/ACM-ICPC-api-service
cd ${ACM_ICPC_api_service} && go get ./...
```
4. 运行
```bash
cd ${ACM-ICPC-api-service} && go build main.go && ./main --config conf.json
```
