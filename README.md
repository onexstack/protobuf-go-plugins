# protobuf-go-plugins

onexstack 自维护的 protoc 插件集合。每个插件都是从上游 fork 出来后可以自己改的副本，
而不是一个 pin 住的上游版本——这样上游的缺陷能在本仓修掉，需要的功能也能直接加。

## 插件

| 插件 | 上游 | 本仓改动 |
|---|---|---|
| `cmd/protoc-gen-go-deepcopy` | [protobuf-tools/protoc-gen-deepcopy](https://github.com/protobuf-tools/protoc-gen-deepcopy) | 修掉生成代码复制 `sync.Mutex` 的缺陷，见下 |
| `cmd/protoc-gen-go-defaults` | [linka-cloud/protoc-gen-defaults](https://github.com/linka-cloud/protoc-gen-defaults) | 修掉 `imports` 跨文件泄漏，见下 |
| `cmd/protoc-gen-go-errors-code` | onexstack/onex `tools/` | 修正悬空 import；目录改名以匹配其标志 |
| `cmd/protoc-gen-go-errors` | [go-kratos/kratos](https://github.com/go-kratos/kratos) | 尚未迁移，自带独立 `go.mod` |
| `cmd/protoc-gen-go-json` | — | 占位，尚无实现 |

顶层还有 `defaults/`：它是 `protoc-gen-go-defaults` 的**运行时**那一半，不是插件。
生成代码会空白导入它，所以它必须放在一个描述其真实身份的路径下，而不是 `cmd/` 里。
详见该插件的 README。

### 迁移时反复出现的同一个坑

三个插件里有三个是「从别处拷进来、import 没跟着改」。这类错误不会在只构建你手上的那个
插件时暴露——只有整个模块被加载时才会，而 `go mod tidy` 正是这么做的：

| 插件 | 遗留的 import | 后果 |
|---|---|---|
| `protoc-gen-go-errordoc` | `onexstack/onex/tools/protoc-gen-go-errors-code/errors`（该路径不存在） | 模块根 `go mod tidy` 失败 |
| `protoc-gen-go-defaults` | 旧 module 路径写进了 `go_package` | 重生成会把 import 写回旧模块 |

判断遗留 import 有个可靠的信号：**它指向的路径在本机根本不存在**。`go build ./...`
在一个可构建的模块里应当零错误；如果它只对某个子目录报错，先怀疑那个子目录的 import。

### protoc-gen-go-deepcopy

为 `.pb.go` 类型生成 `DeepCopyInto()` / `DeepCopy()` / `DeepCopyInterface()`。

上游模板把 `DeepCopyInto` 写成：

```go
func (in *T) DeepCopyInto(out *T) {
	p := proto.Clone(in).(*T)
	*out = *p        // ← 问题在这里
}
```

`*out = *p` 是**结构体整体赋值**，而生成的 message 内嵌了
`protoimpl.MessageState`——其中含一个 `sync.Mutex` 和一个原子指针。整体赋值会把这份
状态一并拷走，于是副本与原对象**共用同一把锁**。这不是 `go vet` 的误报，是任一方的
消息被并发使用时就会发生的数据竞争。

`proto.Clone` 本身没有问题：它用 `New()` + `mergeMessage` 构造结果，从不复制 state。
错的只是跟着的那一行赋值。

本仓改为走 protobuf 反射层：

```go
func (in *T) DeepCopyInto(out *T) {
	proto.Reset(out)
	proto.Merge(out, in)
}
```

`Merge` 逐字段通过反射 API 写入，不碰 state；`Reset` 先清空，使 `out` 非零时不残留
源对象没有的字段。生成文件的 import 与三个方法名均未变，因此消费方无需改动。

## 使用

插件目录名决定 protoc 标志：目录叫 `protoc-gen-X`，protoc 就以 `--X_out` 调用它。
例如 `protoc-gen-go-deepcopy` 对应 `--go-deepcopy_out`。**改目录名等于改标志名**，
消费方必须同步修改。

```bash
make install     # go install 到 $(go env GOPATH)/bin
make build vet test
make help
```

## 已知问题

- `cmd/protoc-gen-go-errors` 自带独立 `go.mod`（沿用上游 `go-kratos` 的 module 路径），
  因此本仓目前同时存在「并入根模块」和「自带 go.mod」两种组织方式。其余插件都已并入根模块。
- `cmd/protoc-gen-go-json` 是空占位，尚无实现。
