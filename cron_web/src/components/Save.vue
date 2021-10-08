<template>
  <div class="box" >
    <div class="save">
      <table id="save" border="2px" cellspacing="0">
        <thead>
        <tr>
          <th>任务名称</th>
          <th>shell命令</th>
          <th>cron表达式</th>
        </tr>
        </thead>

        <tbody>
        <tr>
          <td align="center"><input type="text" style="width: 80%" v-model="task.task_name"></td>
          <td align="center"><input type="text" style="width: 80%" v-model="task.task_command"></td>
          <td align="center"><input type="text" style="width: 80%" v-model="task.task_cron_expr"></td>
        </tr>
        </tbody>
      </table>
    </div>
    <div class="commit">
      <button style="color: darkorange" @click="saveTask"><h2>提交</h2></button>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "Save",
  unmounted() {
    this.$emit("flush")
  },
  data(){
    return {
      task: {
        task_name: "sayHello",
        task_command: "echo hello",
        task_cron_expr: "*/5 * * * * *",
      }
    }
  },
  methods:{
    saveTask(){
      let that = this;
      axios({
        url: "http://192.168.162.128:9090/api/v1/save",
        method: "POST",

        // body
        data: that.task

      }).then(
        res => {
          alert(res.data.message)
          this.$router.go(0)
        },
        err => {
          alert(err.data)
          this.$router.go(0)
        })
    }
  },
}
</script>

<style scoped>
.box{
  background-color: lightblue;
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
}

.save{
  height: 70%;
}

.save #save{
  position: fixed;
  top: 40%;
  margin: auto;
  width: 100%;
}

.commit{
  display: flex;
  position: fixed;
  left: 50%;
}
</style>