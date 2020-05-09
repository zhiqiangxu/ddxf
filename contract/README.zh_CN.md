# ddxf合约

## 核心

链上记录发布、下单、核销动作，链下发货，有纠纷走仲裁

## 流程

1. 卖家通过`DTokenSellerPublish`发布**商品**，可以是一件或者多件**物品**
    
    1. 需要注意的是，每件物品都必须有一个**可验证**的`json`描述哈希(`sha256`，可通过`ddxf`提供的`Sha256Bytes`将任意字节流转为符合要求的描述哈希)
        1. 对于描述中的`map`或者`object`结构，序列化时需要按`key`升序排列, 可通过`ddxf`提供的`Object2Bytes`将任意描述对象转为符合要求的`json`)

    2. 对于静态资源，还需要额外`append`一个**可验证**的数据哈希(`crc32`, 可通过`ddxf`提供的`Crc32Bytes`将任意字节流转为符合要求的数据哈希)


2. 买家通过`BuyDToken`购买卖家的**商品**，或者通过`BuyDTokenFromReseller`从已购买但未使用的买家购买

3. 买家购买后，卖家通过`UseToken`对DToken中的单件**物品**进行核销，并在链外实现买家的权益
    1. 需要注意的是，**链上核销**与**链外实现权益**，天然是非原子的，也是纠纷的源头，所以买卖双方需要各自留好证据。

4. 对于加工服务的场景，买家可以通过`SetDTokenAgents`或者`AddDTokenAgents`设置**商品**的`agent`，这样`agent`便可以通过`UseDTokenByAgent`代替买家消费了；买家也可以通过`RemoveDTokenAgents`取消`agent`的代理资格。