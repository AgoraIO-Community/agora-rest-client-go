## CloudRecording Example

> 这是 Agora Cloud Recording 的一个示例项目，使用了 Agora Cloud Recording RESTful API，实现了频道录制的功能。本示例支持合流录制、单流录制和页面录制三种模式。

### 运行示例项目

#### 前提条件

在当前目录创建一个 `.env` 文件，填入以下内容：

```bash
APP_ID=<Your App ID>
BASIC_AUTH_USERNAME=<Your Basic Auth Username>
BASIC_AUTH_PASSWORD=<Your Basic Auth Password>
CNAME=<Your Channel Name>
TOKEN=<Your Token>
UID=<Your UID>
STORAGE_CONFIG_VENDOR=<Your Storage Vendor>
STORAGE_CONFIG_REGION=<Your Storage Region>
STORAGE_CONFIG_BUCKET=<Your Storage Bucket>
STORAGE_CONFIG_ACCESS_KEY=<Your Storage Access Key>
STORAGE_CONFIG_SECRET_KEY=<Your Storage Secret Key>
```
相关的参数可以通过可在 [CloudRecording 服务说明](../../services/cloudrecording/README.md) 查看

#### 执行

通过下面的命令来运行示例项目：

```bash
go run main.go -mode=<mode>
```

其中 `<mode>` 表示运行的云录制模式，以下是 mode 数值对应的录制模式：
* 1: MixRecording 合流录制
* 2: IndividualRecording 单流录制
* 3: WebRecording 页面录制
