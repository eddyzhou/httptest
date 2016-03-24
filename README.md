About
=====

httptest is a http test dsl that makes http api testing easier.

Command
=====

### set

给内置变量或者自定义变量赋值

内置变量:

- HOST
- POSTDATA
- HEADER

自定义变量格式: **${variable}**

例如：
  
    set HOST http://www.baidu.com
    set POSTDATA username=a,password=b
    set HEADER Qbtoken:token,Uuid:123
    set ${id} 10001
    set ${foo} foo

### req

发送http请求

发送get请求：req GET /xxx

发送post请求：req POST /xxx

### echo

打印变量值

内置变量：

- echo $STATUS

- echo $HOST

- echo $RESP

- echo $POSTDATA

- echo $HEADER

自定义变量：

- echo ${variable}

### getvalue


获取json格式的应答里指定元素的值

例如: 

     http response为{"err": 0, "foo": "foo1", "bar": "bar1"} 
     getvalue ${err} err
     getvalue ${foo} bar
     ${err}为0, ${foo}为"bar1"
     

### getarray

获取json格式的应答里，JSONArray里每个元素指定key的值的数组
 
例如：

    http response为{"err": 0, "users": [{"id":123, "name":"foo"}, {"id": 456, "name": "bar"}]}
    getarray ${ids} id
    ${ids}的值为[123, 456]

### assert

断言变量的值

例如:
     
    assert $STATUS 200
    assert ${foo}  foo
     
### for

循环执行command，以endfor结束，支持循环嵌套

例如：

    for ${typ} in [suggest,txt]
    
    req GET /article/list/${typ}
    getarray ${ids} id
    
    for ${aid} in ${ids}
    req GET /article/${aid}
    endfor
    
    endfor