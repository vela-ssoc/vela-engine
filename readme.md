# engine
> 利用template模式可以统一处理某一类问题的解释器和引擎
> 主要包含脚本模板和解释器

## 内置功能
- [规则模板](#模板) &emsp;规则模板
- [engine](#解释引擎) &emsp;规则模板解释器
- [feedback](#运行结果) &emsp;规则模板运行结果

## 模板
> 模板是是功能核心功能,为实际的扫描引擎提供数据和场景模板,阔以灵活的处理进程,文件, 网络等信息,做批量执行和扫描<br/>
> 模板主要设计参考了[nuclei-templates](https://github.com/projectdiscovery/nuclei-templates) <br />
> 模板格式采用了yaml作为配置文件的方式 主要包含一下信息 <br />
> 字段满足框架中各个对象的index字段信息 <br />

```yaml
id: CVE-2022-001

info:
  name: u5cms v8.3.5 - Open Redirect
  authors: 0x_Akoko,0x_vela
  protocol: process
  severity: medium
  description: |
    u5cms version 8.3.5 contains a URL redirection vulnerability that can cause a user's browser to be redirected to another site via /loginsave.php.
  reference:
    - https://github.com/u5cms/u5cms/issues/50
    - https://nvd.nist.gov/vuln/detail/CVE-2022-32444
  classification:
    cvss-metrics: CVSS:3.1/AV:N/AC:L/PR:N/UI:R/S:C/C:L/I:L/A:N
    cvss-score: 6.1
    cve-id: CVE-2022-32444
    cwe-id: CWE-601

  metadata:
    aaa: 123
    abc: 123

  tags: cve,cve2022,redirect,u5cms,cms,process

condition:
  logic: and
  matchers:
    - type: name
      method: word
      value:
        - java

    - type: pid
      method: ge
      value:
        - 10

before: |-
  vela.keyword("hello")
  
after: |-
  if ctx.feedback then
    ctx.payload("pid:%d exe:%s username:%s" , ctx.value.pid , ctx.value.exe , ctx.value.username)
  end
  
  if ctx.value.username == "root" then
      ctx.happy()
  end
```

## 模板功能介绍
> 模板功能主要包含一下几点需要注意
+ 1.id是命中后的唯一标识
+ 2.protocol是当前模板匹配的对象类型
+ 3.tags可以用来过滤是否需要匹配
+ 4.条件执行包含 regex,equal,word,keyword;其中keyword是当前模块产生的关键字,如regex匹配参数中part
+ 5.before 是在condition之前执行脚本
+ 6.after 在condition执行过执行


### context
模块的周期变量 解决 before 和 after的脚本在执行过程中的一些接口函数

内置字段:

- id           &emsp;模板规则ID值
- file         &emsp;模板配置文件
- description  &emsp;模板配置文件描述
- raw_payload  &emsp;payload的字符串信息
- raw_keyword  &emsp;keyword的字符串信息
- raw_tags     &emsp;规则的标签信息
- data_type    &emsp;数据的类型
- value        &emsp;数据值
- feedback     &emsp;一般存在于after脚本中,用来判断是否命中规则

内置函数:
- ctx.payload(format , ...) &emsp;添加过程匹配的载荷数据
- ctx.happy()               &emsp;匹配模块确定命中
- ctx.match(v)              &emsp;匹配模块内部调用match函数


## 解释引擎
> 加载模块的释器引擎,专门用来批量执行template模块 <br />
>  ee = vela.engine(cfg)

配置字段:

- rules: &emsp; [模板集合](#模板数据)规则文件
- tags: &emsp;  根据标签值来过滤加载的配置文件,可以利用*最模式匹配

内置方法:

- ee.case(cond)   &emsp; 执行过程[feedback](#feedback)字段触发条件
- ee.scan(object) &emsp;扫描某个对象
- ee.with([feedback](#feedback)) &emsp;根据配置文件自动扫描无须外部传递扫描对象

newindex:

- rules: &emsp;修改模块的&nbsp;[ee.rules={"xxx" , "bbb"}](#)
- tags : &emsp;修改模块的&nbsp;[ee.tags={"xxx" , "bbb*"}](#)

样例脚本：
```lua
    local e = vela.engine{
        rules = "rule.d/CVE-2022-*.yaml",
        tags  = "process"
    }

    local fb = vela.engine.feedback()
    e.case("feedback").pipe(fb.collect) --feedback = true 触发 condition

    -- 获取process
    vela.ps().pipe("*" , s.scan)

```

## 运行结果
> 模块执行完后返回的数据信息

内部字段:
- [collect(object)](#)  &emsp; 收集[template.context](#context)信息
- [debug()](#)          &emsp; 打印当前所有数据调试结果
- [pipe(context)](#)    &emsp; 遍历所有结果[template.context](#context)并处理


## 样例1.扫描传入对象
> 直接传入扫描对象值 利用scan的方法执行

```lua
    local e = vela.engine{
        rules = "rule.d/CVE-2022-*.yaml",
        tags  = "process"
    }
    local fb = vela.engine.feedback()
    e.case("feedback").pipe(fb.collect)

    -- 获取process
    vela.ps().pipe("*" , s.scan)

    fb.pipe(function(ctx)  --处理返回扫描的返回结果
        
    end)
```


## 样例2.直接启动规则模板
> 利用配置文件自动 扫描
```lua
    local e = vela.engine{
        rules = "rule.d/CVE-2022-*.yaml",
        tags  = "process"
    }
    local fb = vela.engine.feedback()
    e.with(fb) --直接启动解释器, 因为有些before脚本可以获取自己想要获取的信息

    fb.pipe(function(ctx)
        --todo
    end)
```
