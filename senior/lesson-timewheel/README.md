# 时间轮算法实现定时器

## 原理
比如：2023-09-21 08:30:00，实际上，就是分成了 `{year}-{month}-{date}-{hour}-{minute}-{second}` 组成的 6 级时间轮等级结构。

## 文档
https://zhuanlan.zhihu.com/p/658079556

## 参考资料
https://github.com/HDT3213/godis
https://github.com/xiaoxuxiansheng/xtimer
https://github.com/xiaoxuxiansheng/timewheel/blob/main/time_wheel_test.go