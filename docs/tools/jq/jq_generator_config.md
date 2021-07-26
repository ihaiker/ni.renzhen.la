# JQ 工具快速生成JSON文件

>   有些时候我们需要快速生成一个Json文件，此JSON文件的部分内容来自env

**json.template.json** 文件

```json
{
  "region": env.REGSON,
  "account": env.ACCOUNT
}
```

我们执行下面的shell:

```bash
REGION=cn-hangzhou ACCOUNT=demo jq -n -f json.template.json 
```

后会直接修改编辑并输出结果文件:

```json
{
  "region": "cn-zhangzhou",
  "account": "demo"
}
```

此文件就可以直接用于你的实例中。