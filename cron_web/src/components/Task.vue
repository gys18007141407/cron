<template>
<div class="task">
  <table id="task" v-if="!isEmpty" border="2px" cellspacing="0">
    <thead>
    <tr>
      <th>任务名称</th>
      <th>shell命令</th>
      <th>cron表达式</th>
      <th>任务控制</th>
    </tr>
    </thead>

    <tbody>
    <tr v-for="(task, idx) in this.data">
      <td><input type="text" v-model="task.task_name" style="width: 80%" v-bind:disabled="!editable[idx]" ></td>
      <td><input type="text" v-model="task.task_command" style="width: 80%" v-bind:disabled="!editable[idx]"></td>
      <td><input type="text" v-model="task.task_cron_expr" style="width: 80%" v-bind:disabled="!editable[idx]"></td>
      <td>
        <div class="ctrl">
          <button class="btn" @click="editTask(idx)" v-show="!editable[idx]">编辑</button>
          <button class="btn" @click="saveTask(idx)" v-show="editable[idx]">提交</button>
          <button class="btn" @click="taskLog(idx)">日志</button>
          <button class="btn" @click="killTask(idx)" style="color: lightcoral">杀死</button>
          <button class="btn" @click="removeTask(idx)" style="color: red">移除</button>
        </div>
      </td>
    </tr>
    </tbody>
  </table>
  <div v-else>没有任务，是否<button @click="newTask">新建任务</button></div>
</div>
</template>

<script>
import axios from "axios";

export default {
  name: "Task",
  props:{
    data:{
      type: Array,
      required: true
    },
    message:{
      type: String,
      required: true
    },
    errno:{
      type: Number,
      required: true
    }
  },
  data(){
    return{
      editable:[]
    }
  },

  computed:{
    isEmpty(){
      return this.data.length === 0
    },
  },
  methods:{
    editTask(idx){
      this.editable[idx] = true
    },
    saveTask(idx){
      this.editable[idx] = false
      let that = this
      axios({
        url: "http://192.168.162.128:9090/api/v1/save",
        method: "POST",

        // body
        data: that.data[idx]

      }).then(
        res => {
          alert(res.data.message)
          this.$router.go(0)
        },
        err => {
          alert(err.data)
          this.$router.go(0)
        })
    },
    taskLog(idx){
      this.$emit("getLog", this.data[idx].task_name)
    },
    killTask(idx){
      let that = this
      axios({
        url: "http://192.168.162.128:9090/api/v1/kill",
        method: "POST",

        // header
        // headers:{
        //   task_name: that.data[idx].task_name,
        // },

        // query
        params:{
          task_name: that.data[idx].task_name,
        },

        // body
        data: {
          task_name: that.data[idx].task_name,
        }
      }).then(
        res => {
          alert(res.data.message)
          this.$router.go(0)
        },
        err => {
          alert(err.data)
          this.$router.go(0)
        })
    },
    removeTask(idx){
      let that = this
      axios({
        url: "http://192.168.162.128:9090/api/v1/delete",
        method: "DELETE",

        // query
        params:{
          task_name: that.data[idx].task_name,
        },

      }).then(
        res => {
          alert(res.data.message)
          this.$router.go(0)
        },
        err => {
          alert(err.data)
          this.$router.go(0)
        })
    },
    newTask(){
      this.$router.push('/home/save')
    },
  }
}
</script>

<style scoped>
.task{
  background-color: wheat;
  text-align: center;
  height: 100%;

}

#task{
  width: 100%;
  table-layout: fixed;
}

#task tr td{
}

.ctrl{
  height: 100%;
}

.btn{
  color: deepskyblue;
  position: center;
}

</style>