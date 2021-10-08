<template>
  <div class="home" v-show="!save">
    <div class="banner">分布式定时任务管理</div>
    <div class="content">
      <div class="menu">
        <div class="nav">
          <button class="btn" @click="listTask">任务列表</button>
        </div>
        <div class="nav">
          <button class="btn" @click="listWorker">健康节点</button>
        </div>
        <div class="nav">
          <button class="btn" @click="newTask">新建任务</button>
        </div>
      </div>
      <div class="message">
        <router-view v-bind:data="data" :message="message" :errno="errno" @getLog="listLogSwitchRouter" @listLog="listLog" @flush="listTask"></router-view>
      </div>
    </div>
  </div>

</template>

<script>

import axios from "axios"

export default {
  name: 'Home',
  data(){
    return{
      data: [],
      errno: 0,
      message: "",

      task_name: "",
    }
  },

  created() {
    this.save = false;
    this.listTask()
  },

  methods:{
    listTask(){
      let that = this
      axios.get("http://192.168.162.128:9090/api/v1/list").then(
          function (response){
            that.data = response.data.data
            that.errno = response.data.errno
            that.message = response.data.message
            that.$router.push('/home/tasks')
          },
          function (err){
            alert(err)
          }
      )
    },
    listWorker(){
      let that = this
      axios.get("http://192.168.162.128:9090/api/v1/worker").then(
          function (response){
            that.data = response.data.data
            that.errno = response.data.errno
            that.message = response.data.message
            that.$router.push('/home/workers')
          },
          function (err){
            alert(err)
          }
      )
    },
    newTask(){
      this.$router.push('/home/save')
    },
    listLogSwitchRouter(task_name){
      this.task_name = task_name
      let that = this
      axios({
        url: "http://192.168.162.128:9090/api/v1/log",
        method: "GET",
        params:{
          task_name: that.task_name,
          skip: 0,
          limit: 10,
        }
      }).then(
          res => {
            that.data = res.data.data
            that.errno = res.data.errno
            that.message = res.data.message
            that.$router.push('/home/log')
          },
          err => {
            alert(err)
          }
      )
    },
    listLog(idx){
      let that = this
      axios({
        url: "http://192.168.162.128:9090/api/v1/log",
        method: "GET",
        params:{
          task_name: that.task_name,
          skip: idx*10,
          limit: 10,
        }
      }).then(
          res => {
            that.data = res.data.data
            that.errno = res.data.errno
            that.message = res.data.message
          },
          err => {
            alert(err)
          }
      )
    }
  }
}
</script>

<style scoped>
.banner{
  background-color: cadetblue;
  text-align: center;
  line-height: 50px;
  color: black;
  height: 50px;
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
}

.content{
  height: 600px;
  position: fixed;
  left: 0;
  right: 0;
  top: 50px;
  display: flex;
}
.menu{
  position: fixed;
  left: 0;
  right: 90%;
  top: 50px;
  background-color: darkgrey;
  height: 600px;
}

.message{
  position: fixed;
  left: 10%;
  right: 0;
  top: 50px;
  background-color: lightskyblue;
  height: 600px;
}

.nav{
  box-sizing: border-box;
  text-align: center;
  width: 100%;
  background-color: beige;
}

.btn{
  width: 100%;
  height: 30px;
  color: blue;
}

</style>