# 中英文字符的字节数占用（byte数组查看即可）

Go 语言的字符串都以 UTF-8 格式保存，每个中文占用 3 个字节，因此使用 len() 获得两个中文文字对应的 6 个字节。

* len(字符串)：查看的是ASCII码数量，一个ASCII码就是一个字节，所以相当于查看的字节数。
* utf8.RuneCountInString(字符串)：Go 语言中 UTF-8 包提供的 RuneCountInString() 函数，统计 Uncode 字符数量。