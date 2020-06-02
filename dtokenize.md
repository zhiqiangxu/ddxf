# dtoken a-z

假设有2个`Data`，`data`和`data'`：

```
{
    "数据属性1":"属性值1",
    "数据属性2":"属性值2",
    "数据属性3":"属性值3",
    "Endpoint":"endpoint1",
}
```

和


```
{
    "数据属性1'":"属性值1'",
    "数据属性2'":"属性值2'",
    "数据属性3'":"属性值3'",
    "Endpoint":"endpoint2",
}
```

为每个`Data`设计2个`Token`

`data`的`Token`为：

```
{
    "权限1":"权限值1",
    "Endpoint":"endpoint3",
}
```

和

```
{
    "权限2":"权限值2",
    "Endpoint":"endpoint4",
}
```

`data'`的`Token`为：

```
{
    "权限1'":"权限值1'",
    "Endpoint":"endpoint5",
}
```

和

```
{
    "权限2'":"权限值2'",
    "Endpoint":"endpoint6",
}
```

假设`MP Item`（MP露出的商品信息）的格式是：


```
{
    "Item属性1":"Item值1",
    "Item属性2":"Item值2",
}
```

假设将之前准备的`Data`和`Token`都选中，也就是希望出售的是2个`Data`的所有权限，对应：

```
[
    {
        "data":{
            "数据属性1":"属性值1",
            "数据属性2":"属性值2",
            "数据属性3":"属性值3",
            "Endpoint":"endpoint1",

        },
        "token":[
            {
                "权限1":"权限值1",
                "Endpoint":"endpoint3",
            },
            {
                "权限2":"权限值2",
                "Endpoint":"endpoint4",
            }
        ]
    },
    {
        "data":{
            "数据属性1'":"属性值1'",
            "数据属性2'":"属性值2'",
            "数据属性3'":"属性值3'",
            "Endpoint":"endpoint2",
        },
        "token":[
            {
                "权限1'":"权限值1'",
                "Endpoint":"endpoint5",
            },
            {
                "权限2'":"权限值2'",
                "Endpoint":"endpoint6",
            }
        ]
    }
]
```




上面的json用于生成链上的token，其实就是`token template`。`Item+token template`，加上Fee、ExpiredDate、Stocks等售卖信息后，即为链上信息，包括了整个过程的所有信息，此外为了减少上链数据量，实际是把`Item`、`Data`和`Token`的可验证哈希以及Endpoint上链，通过哈希可以从对应的Endpoint查出原始数据。传统互联网中占据核心地位的MP，此时仅包含上面的Item信息，用于商品的露出。
