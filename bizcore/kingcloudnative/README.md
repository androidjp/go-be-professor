# King维度的云原生能力搭建

能力包括：（主要适应king企业）
* 服务cicd相关配置
* 动态配置能力（apollo、etcd、环境变量、配置文件）
* 链路追踪
* prom上报
* dw上报
* 服务发现与grpc调用能力【TODO】

## 目录说明
* neometric: 企业上报埋点收集
* prommetric: prometheus自定义上报
* scfg: 用的支持多渠道的配置读取
* slogger: 基于g1.21版本 slog库的简单封装，支持配置化context字段读取和打印。
* tracing: 全流程链路追踪相关
