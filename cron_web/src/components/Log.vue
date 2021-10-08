<template>
<div class="log">
  <table id="log" v-if="!isEmpty" border="2px" cellspacing="0">
    <thead>
    <tr>
      <th>任务名称</th>
      <th>shell命令</th>
      <th>任务标准输出</th>
      <th>任务错误输出</th>
      <th>计划调度时间</th>
      <th>实际调度时间</th>
      <th>开始运行时间</th>
      <th>结束运行时间</th>
    </tr>
    </thead>

    <tbody>
    <tr v-for="(log, idx) in this.data">
      <td align="center">{{log.task_name}}</td>
      <td align="center">{{log.task_command}}</td>
      <td align="center">{{log.task_output.length ? log.task_output : "空"}}</td>
      <td align="center">{{log.task_error.length ? log.task_error : "空"}}</td>
      <td align="center">{{this.getTime(log.schedule_time)}}</td>
      <td align="center">{{this.getTime(log.real_schedule_time)}}</td>
      <td align="center">{{this.getTime(log.exec_time)}}</td>
      <td align="center">{{this.getTime(log.finish_time)}}</td>
    </tr>
    </tbody>
  </table>
  <div v-else><b>没有日志</b></div>

  <div v-show="this.data.length > 0" style="background-color: lightgrey">
    <button @click="log0">首页</button>
    <button @click="logPre">上一页</button>
    <input type="text" :value="cur+1" style="text-align: center; width: 10px; margin: 0">
    <button @click="logNext">下一页</button>
  </div>

</div>
</template>

<script>
import axios from "axios";

export default {
  name: "Log",
  props:{
    data:{
      type: Array,
      required: true
    }
  },
  data(){
    return {
      task_name: "",
      cur: 0,
    }
  },
  computed:{
    isEmpty(){
      return this.data.length === 0;
    }
  },

  unmounted() {
    this.$emit("flush")
  },

  methods:{
    getTime(millSec){
      let data = new Date(millSec)
      let y = data.getFullYear()
      let m = data.getMonth()+1
      let d = data.getDate()
      let h = data.getHours()
      let min = data.getMinutes()
      let s = data.getSeconds()
      let ms = data.getMilliseconds()
      return y+"."+m+'.'+d+'-['+h+':'+min+':'+s+'.'+ms+']'
    },
    getLog(page){
      let that = this;
      console.log(page)
      axios({
        url: "http://192.168.162.128:9090/api/v1/log",
        method: "GET",
        params:{
          task_name: that.task_name,
          skip: page*10,
          limit: 10,
        }
      }).then(
          res => {
            if(res.data.length === 0){
              alert("已经是最后一页了")
            }else {
              that.data = res.data
              that.cur = page+1
            }
          },
          err => {
            alert(err)
          }
      )
    },
    log0(){
      this.$emit("listLog", 0)
      this.cur = 0
    },
    logNext(){
      if(this.data.length < 10) alert("已经是最后一页了")
      else {
        this.$emit("listLog", this.cur + 1)
        this.cur ++
      }
    },
    logPre(){
      if(this.cur === 0) alert("已经是第一页了")
      else {
        this.$emit("listLog", this.cur-1)
        this.cur --
      }
    }
  }
}
</script>

<style scoped>
.log{
  background-color: #f6f6f6;
  text-align: center;
}

#log{
  width: 100%;
}
</style>