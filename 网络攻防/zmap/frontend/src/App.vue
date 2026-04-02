<template>
  <div class="app">
    <div class="container">
      <label for="target-ip" class="label">目标IP</label>
      <input v-model="data.ip" id="target-ip" class="input" placeholder="请输入目标IP或范围" />

      <div class="button-group">
        <button
            :class="['button', { 'active': activeipButton === 'ipmethon1' }]"
            @click="setIpMethod('1')">单一主机
        </button>
        <button
            :class="['button', { 'active': activeipButton === 'ipmethon2' }]"
            @click="setIpMethod('2')">目标范围
        </button>
      </div>
      <hr>
      <div>
        <label for="port-value" class="label">请输入端口值或范围</label>
        <input v-model="data.port" :disabled="flag" id="port-value" class="input" />
      </div>
      <div class="button-group">
        <button
            :class="['button', { 'active': activeportButton === 'portmethod1' }]"
            @click="setPortMethod('1');data.port=''">全端口(默认模式)
        </button>
        <!--      <button-->
        <!--          :class="['button', { 'active': activeportButton === 'portmethod2' }]"-->
        <!--          @click="setPortMethod('2')">范围端口-->
        <!--      </button>-->
        <button
            :class="['button', { 'active': activeportButton === 'portmethod3' }]"
            @click="setPortMethod('3')">单一端口
        </button>
      </div>
      <hr>
      <label class="label">确认扫描方式</label>
      <div class="button-group">
        <button
            :class="['button', { 'active': activescanButton === 'scanmethod1' }]"
            @click="setScanMethod('1')">全连接
        </button>
        <button
            :class="['button', { 'active': activescanButton === 'scanmethod2' }]"
            @click="setScanMethod('2')">SYN
        </button>
        <button
            :class="['button', { 'active': activescanButton === 'scanmethod3' }]"
            @click="setScanMethod('3')">FIN
        </button>
        <!--      <button-->
        <!--          :class="['button', { 'active': activescanButton === 'scanmethod4' }]"-->
        <!--          @click="setScanMethod('4')">NULL-->
        <!--      </button>-->
        <!--      <button-->
        <!--          :class="['button', { 'active': activescanButton === 'scanmethod5' }]"-->
        <!--          @click="setScanMethod('5')">XMAS-->
        <!--      </button>-->
        <!--      <button-->
        <!--          :class="['button', { 'active': activescanButton === 'scanmethod6' }]"-->
        <!--          @click="setScanMethod('6')">ACK-->
        <!--      </button>-->
      </div>
      <hr>
      <div>
        <label for="scan-count" class="label">扫描次数:</label>
        <input v-model="data.scannum" id="scan-count" class="input" placeholder="默认为一次" />
      </div>
      <button class="button" @click="test">开始扫描</button>
    </div>
    <div class="display">
      <label class="label">扫描结果:</label>
      <textarea v-model="result" class="dis" readonly></textarea>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {compile, reactive, ref} from "vue";
import axios, {toFormData} from "axios";

// 定义响应式变量
const result=ref('')
const flag=ref(true)
const activeipButton = ref('');
const activeportButton = ref('');
const activescanButton = ref('');
const data = reactive({
  ip: '',
  port: '',
  scannum: '',
  ipmethon: '',
  portmethon: '1',
  scanmethon: ''
})

function setIpMethod(method: string) {
  data.ipmethon = method;
  activeipButton.value = method === '1' ? 'ipmethon1' : 'ipmethon2';
}

function setPortMethod(method: string) {
  data.portmethon = method;
  activeportButton.value = method === '1' ? 'portmethod1' : method === '2' ? 'portmethod2' : 'portmethod3';

  if (method === '1') {
    flag.value = true;  // 禁用
  } else {
    flag.value = false;  // 允许输入
  }
}

function setScanMethod(method: string) {
  data.scanmethon = method;
  activescanButton.value = `scanmethod${method}`;
}

// 点击开始扫描按钮时的逻辑
function test() {
  console.log(data);
  if(data.ip==''){
    alert("ip未填")
    return
  }
  if(data.port==''&&data.portmethon!=='1'){
    alert("端口未填")
    return
  }
  if(data.ipmethon===''){
    alert("目标类型未选")
    return
  }
  if(data.scanmethon===''){
    alert("扫描方式未选")
    return
  }
  if (data.ipmethon === '1') {
    const partan = /^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;

    // if (!partan.test(data.ip)) {
    //   alert("请输入正确的IPv4地址");
    //   return
    // }
  }else {
    const ipv4WithCIDR = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\/(3[0-2]|[12]?[0-9])$/;
    if(!ipv4WithCIDR.test(data.ip)){
      alert("错误地址范围")
      return
    }
  }
  if(data.portmethon === '2'){
    const portparten=/^(1|[1-5][0-9]{0,4}|[1-5][0-9]{0,4})-(6553[0-5]|655[0-2][0-9]|654[0-9]{2}|65[0-4][0-9]{2}|6[0-4][0-9]{3}|[1-5][0-9]{0,4}|[1-9][0-9]{0,4})$/;
    if(!portparten.test(data.port)){
      alert("端口范围错误")
      return
    }
  }else if (data.portmethon === '3') {
    const ports = data.port.split(',').map(Number);
    for (const value of ports) { // 使用 `of` 来遍历数组
      if (value <= 0 || value > 65535) {
        alert("含有错误端口");
        return;
      }
    }
  }
  // 确保扫描次数是一个非负整数
  const scanCountPattern = /^[1-9]\d*$/; // 大于0的正整数
  if (!scanCountPattern.test(data.scannum) && data.scannum !=='') {
    alert("非法的次数");
    return;

  }
  data.scannum='1'
  const options = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json; charset=utf-8',
    },
    url: 'http://localhost:7899', // 确保这个是有效的字符串
    data: data,
  };

  axios(options).then(res=>{
    console.log(res);
    // 假设返回的数据是上面提供的 JSON 对象
    const re= res.data; // 获取返回的数据

    // 格式化 JSON 对象为字符串
    let formattedOutput = '';
    for (const [key, value] of Object.entries(re)) {
      formattedOutput += `${key}: ${value}`; // 格式化为 "IP: 状态"
    }

    // 将格式化后的字符串赋值给 output
    result.value = formattedOutput;
    console.log(result)
  })
}
</script>

<style>
.app{
  display: flex;
  margin-top: 10%;
  width: 100%; /* 容器宽度 */
  max-width: 600px; /* 最大宽度 */
  padding: 30px; /* 内边距 */
  /* 上下左右边距 */
  background-color: #ffffff; /* 容器背景色 */
  border-radius: 10px; /* 圆角 */
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2); /* 阴影效果 */
}
body {
  margin: 0; /* 去除默认边距 */
  padding: 0; /* 去除默认内边距 */
  background-color: #e0e0e0; /* 整体背景色 */
  height: 100vh; /* 高度为视口高度 */
  display: flex; /* 使用Flexbox布局 */
  justify-content: center; /* 水平居中 */
  align-items: center; /* 垂直居中 */
  font-family: 'Arial', sans-serif; /* 字体 */
}

.container {
    width: 100%;
}

.label {
  display: block; /* 块级元素 */
  margin-bottom: 5px; /* 下边距 */
  font-weight: bold; /* 加粗字体 */
  color: #1b2636;
}

.button.active {
  background-color: #59ac59; /* 点击后变成绿色 */
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2); /* 悬停时阴影效果 */
}

.input {
  width: 90%; /* 输入框宽度 */
  text-align: center;
  padding: 12px; /* 内边距 */
  margin-bottom: 15px; /* 下边距 */
  border: 1px solid #ccc; /* 边框 */
  border-radius: 5px; /* 圆角 */
  background-color: #f9f9f9; /* 输入框背景色 */
  font-size: 16px; /* 字体大小 */
}

.button-group {
  margin-bottom: 15px; /* 下边距 */
}

.button {
  padding: 10px 15px; /* 内边距 */
  margin-right: 5px; /* 右边距 */
  border: none; /* 去除边框 */
  border-radius: 5px; /* 圆角 */
  background-color: #007bff; /* 按钮背景色 */
  color: #ffffff; /* 字体颜色 */
  cursor: pointer; /* 鼠标样式 */
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2); /* 阴影效果 */
  transition: background-color 0.3s, box-shadow 0.3s; /* 过渡效果 */
  font-size: 16px; /* 字体大小 */
}

.button:hover {
  //background-color: #0056b3; /* 悬停时背景色 */
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.3); /* 悬停时阴影效果 */
}
.display {
  margin-left: 50px;
}

.dis {
  height: 60vh;
  width: 250px;
  font-size: larger;
  font-family: Consolas;
  color: #59ac59;
  text-align: left; /* 确保文本靠左 */
  padding: 10px; /* 添加内边距，避免文本紧贴边缘 */
  line-height: 1.5; /* 增加行距，便于阅读 */
  overflow-y: auto; /* 使内容超过高度时出现滚动条 */
  white-space: pre-wrap; /* 保留换行和空格 */
  background-color: #f1f1f1; /* 设置背景色，提升可读性 */
  border-radius: 5px; /* 添加圆角 */
  resize: none; /* 禁止用户调整大小 */
}

</style>
