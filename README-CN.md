# gtr 

[English](README.md) | [中文](README-CN.md)

gtr 是一个用Go语言编写的终端工具集，提供快速的数据转换和格式互转功能。

## 安装

```bash
go build -o gtr
```

## 使用说明

### 1. 坐标转换

在不同坐标系之间进行转换：WGS84、GCJ02（火星坐标系）和BD09（百度坐标系）。

**语法：**
```bash
gtr coordinate [type] <经度,纬度>
gtr coordinate [type] <经度|纬度>
gtr coordinate [type] <经度> <纬度>
```

**坐标系类型：** `wgs`, `wgs84`, `gcj`, `gcj02`, `gd`, `bd`, `bd09`, `gg`

**使用例子：**
```bash
gtr coordinate 113.901495,22.499501
gtr coordinate gcj 113.901495,22.499501
gtr coordinate wgs 113.901495 22.499501
```

**输出：** 显示转换到三个坐标系统的结果，精度到小数点后6位。

---

### 2. HTTP请求

发送GET或POST HTTP请求，查看详细的响应信息。

**语法：**
```bash
gtr http <URL>
gtr http <URL> '<JSON数据>'
```

**使用例子：**
```bash
gtr http https://api.example.com/endpoint
gtr http https://api.example.com/post '{"key":"value","name":"test"}'
```

**输出：** 显示请求URL、POST数据、响应头、耗时、状态码和响应体。

---

### 3. 时间转换

在时间戳（10位和13位）和人类可读的日期格式之间进行转换。

**语法：**
```bash
gtr time <时间戳|日期>
```

**支持的格式：**
- 10位Unix时间戳：`1727087511`
- 13位毫秒时间戳：`1727087511000`
- 标准格式：`2024-09-23 10:31:51`
- 紧凑格式：`20240923103151`

**使用例子：**
```bash
gtr time 1727087511
gtr time 1727087511000
gtr time "2024-09-23 10:31:51"
gtr time 20240923103151
```

**输出：** 显示所有格式转换结果，包括10位和13位时间戳。

---

### 4. 文本编解码

使用Base64、URL编码和MD5哈希进行文本转换。可以单独执行各种操作，或者用一条命令一次性转换所有格式。

#### 快速转换（所有格式）
```bash
gtr text "<文本>"
```

将显示所有转换结果：
- Base64 编码/解码
- URL 编码/解码  
- MD5 哈希
- 转换失败的操作用 ✗ 标记

**使用例子：**
```bash
gtr text "hello world"
```

#### 单独操作

#### Base64 编码/解码
```bash
gtr text base64 encode "<文本>"
gtr text base64 decode "<base64字符串>"
```

**使用例子：**
```bash
gtr text base64 encode "hello world"
gtr text base64 decode "aGVsbG8gd29ybGQ="
```

#### URL 编码/解码
```bash
gtr text url encode "<文本>"
gtr text url decode "<URL编码字符串>"
```

**使用例子：**
```bash
gtr text url encode "hello world"
gtr text url decode "hello%20world"
```

#### MD5 哈希
```bash
gtr text md5 "<文本>"
```

**使用例子：**
```bash
gtr text md5 "password"
```

---

## 命令别名

所有命令都支持多个别名，方便快速使用：

- `time` (时间转换)：`t`, `-t`, `--t`, `-time`, `--time`
- `coordinate` (坐标转换)：`coor`, `-coor`, `--coor`, `c`, `-c`, `--c`
- `http` (HTTP请求)：`-http`, `--http`, `h`, `-h`, `--h`
- `text` (文本编解码)：`t`, `-t`, `--t`, `-text`, `--text`

**使用例子：**
```bash
gtr t 1727087511                    # 等同于：gtr time 1727087511
gtr coor 113.901495,22.499501       # 等同于：gtr coordinate 113.901495,22.499501
gtr h https://example.com           # 等同于：gtr http https://example.com
```

---

## 功能特性

✅ 支持多种坐标系统转换  
✅ HTTP请求测试，显示详细响应信息  
✅ 灵活的时间戳和日期格式转换  
✅ Base64、URL和MD5文本转换  
✅ 命令别名，快速访问  
✅ 清晰格式化的输出结果

---

## 许可证

请查看 [LICENSE](LICENSE) 文件
