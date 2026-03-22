## Go 语言特性

### 核心设计理念
Go 由 Google 设计，目标是**简洁、高效、并发**，解决大规模工程问题。

### 主要特性

**1. 语法简洁**
- 没有类、继承、泛型（1.18 前）、异常等复杂概念
- 强制代码风格（`gofmt`），减少团队争议
- 编译速度极快

**2. 原生并发支持**
- **Goroutine**：轻量级协程，`go func()` 一行启动，初始栈仅 2KB
- **Channel**：goroutine 间通信，"不要通过共享内存通信，而要通过通信共享内存"
- **select**：多路 channel 复用

**3. 内存管理**
- 自动 GC，但延迟比 JVM 更低、更可预测
- 值类型 + 指针，减少堆分配压力

**4. 接口系统（隐式实现）**
- 不需要 `implements` 关键字，只要实现了接口的方法就自动满足接口
- 面向行为编程，解耦更彻底

**5. 错误处理**
- 用多返回值 `(result, error)` 显式处理错误，而非 try/catch
- 强迫开发者正视每一个错误

**6. 工具链完整**
- 内置测试、性能分析、文档生成、代码格式化

---

## Go vs Java 对比

| 维度 | Go | Java |
|---|---|---|
| **编程范式** | 过程式 + 接口 | 面向对象（类/继承） |
| **并发模型** | Goroutine + Channel（CSP） | Thread / `synchronized` / `CompletableFuture` |
| **类型系统** | 静态类型，接口隐式实现 | 静态类型，显式 `implements` |
| **错误处理** | 多返回值 `(T, error)` | `try/catch/finally` 异常机制 |
| **运行方式** | 编译为原生二进制，无需运行时 | 编译为字节码，需要 JVM |
| **启动速度** | 毫秒级 | 秒级（JVM 预热） |
| **内存占用** | 极低（几 MB） | 较高（JVM 本身占用大） |
| **泛型** | 1.18+ 支持，较简单 | 完整泛型支持 |
| **继承** | ❌ 无继承，用组合替代 | ✅ 单继承 + 接口 |
| **空指针** | 有 nil，需注意 | NullPointerException 是经典痛点 |
| **包管理** | Go Modules（内置） | Maven / Gradle |
| **部署** | 单一可执行文件 | Jar + JVM 环境 |

---

## 代码对比示例

**并发：启动 10 个任务**

```go
// Go - Goroutine
for i := 0; i < 10; i++ {
    go func(id int) {
        fmt.Println("task", id)
    }(i)
}
```

```java
// Java - Thread
for (int i = 0; i < 10; i++) {
    final int id = i;
    new Thread(() -> System.out.println("task " + id)).start();
}
```

**接口实现**

```go
// Go - 隐式实现，无需声明
type Animal interface { Speak() string }

type Dog struct{}
func (d Dog) Speak() string { return "Woof" } // 自动实现 Animal
```

```java
// Java - 显式声明
interface Animal { String speak(); }

class Dog implements Animal {  // 必须写 implements
    public String speak() { return "Woof"; }
}
```

**错误处理**

```go
// Go
data, err := os.ReadFile("file.txt")
if err != nil {
    log.Fatal(err)
}
```

```java
// Java
try {
    byte[] data = Files.readAllBytes(Path.of("file.txt"));
} catch (IOException e) {
    e.printStackTrace();
}
```

---

## 学习建议

作为新手，如果你有 Java 基础，学 Go 时需要重点转变这几个思维：

1. **忘掉类和继承**，改用结构体（struct）+ 组合
2. **接口是鸭子类型**，不用提前声明实现关系
3. **错误不抛出，而是返回**，每个 error 都要处理
4. **并发用 goroutine**，比线程轻太多，大胆用
5. **代码简洁是美德**，Go 刻意移除了很多"高级"特性
