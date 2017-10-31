# What
提供ACM-ICPC赛时服务。服务项：
1. 榜单服务
2. 气球服务
3. 打印机服务
4. 后台管理服务

# Usage

下面是使用方法，具体到代码请参照后面的 Usage Example
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
        "JWTSecret": "terceSTWJ",       \\加密用token，使用一个复杂的长密码串即可
        "Addr": ":8080",                \\端口号，不推荐改变
        "Admin": {
            "Account": "admin",         \\管理员账号
            "Password": "ElPsyCongroo"  \\管理员密码
        },
        "IsTestMode": true,             \\是否为测试模式，测试模式性能较差，用于DEBUG
        "NeedAuth": false,              \\是否需要验证身份，比赛是一定要改成true，重启服务
        "MaxAllowed": 1000              \\最大同时连接数
    },
    "Storage": {
        "Dirver": "sqlite3",            \\数据库名，不得改变
        "Config": "sqlite3.db",         \\sqlite文件名
        "MaxIdleConns": 600,            \\数据库最大空闲连接数
        "MaxOpenConns": 1000            \\数据库最大打开连接数
    },
    "Printer": {
        "QueueSize": 1000,              \\打印机队列长度
        "PrinterNameList": []           \\打印机名，服务器配好打印机后，将打印机名填写到这里即可
    },
    "ResultsXMLPath": "./results.xml",  \\pc2 产出的 results.xml 路径
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
## Fedora 25
0. 以管理员身份运行
```bash
sudo su
```
1. 安装git, golang, sqlite3
```bash
dnf install git golang sqlite3
# 配置gopath，在~/.bashrc中写入GOPATH变量，可参照下面方法
echo "export GOPATH=\"/usr/share/gocode\"" >> ~/.bashrc
source ~/.bashrc
```
2. 下载代码
```bash
WORK_DIR=/home/shjwudp/workstation
mkdir -p ${WORK_DIR}
cd ${WORK_DIR} && git clone git@github.com:shjwudp/ACM-ICPC-api-service.git
```
3. 获取go代码依赖
```bash
ACM_ICPC_api_service=${WORK_DIR}/ACM-ICPC-api-service
cd ${ACM_ICPC_api_service} && go get ./...
```
4. 运行
```bash
cd ${ACM_ICPC_api_service} && sh -x run.sh
```

## Ubuntu 16.04 LTS
0. 以管理员身份运行
```bash
sudo su
```
1. 安装git, golang, sqlite3
```bash
apt-get update
apt-get install git golang sqlite3
# 配置gopath，在~/.bashrc中写入GOPATH变量，可参照下面方法
echo "export GOPATH=\"/usr/share/go\"" >> ~/.bashrc
source ~/.bashrc
```
2. 下载代码
```bash
WORK_DIR=/home/shjwudp/workstation
mkdir -p ${WORK_DIR}
cd ${WORK_DIR} && git clone git@github.com:shjwudp/ACM-ICPC-api-service.git
```
3. 获取go代码依赖
```bash
ACM_ICPC_api_service=${WORK_DIR}/ACM-ICPC-api-service
cd ${ACM_ICPC_api_service} && go get ./...
```
4. 运行
```bash
cd ${ACM_ICPC_api_service} && sh -x run.sh
```