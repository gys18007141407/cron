<template>
  <div class="worker">
    <table id="worker" v-if="!isEmpty" border="2px" cellspacing="0" align="center">
      <thead>
      <tr>
        <th style="color: darkred">节点IP</th>
      </tr>
      </thead>

      <tbody>
      <tr v-for="(worker, idx) in this.data">
        <td>{{worker}}</td>
      </tr>
      </tbody>
    </table>
    <div v-else><b>当前没有任何在线的健康节点</b></div>
  </div>
</template>

<script>
export default {
  name: "Worker",
  props:{
    data:{
      type: Array
    },
    message:{
      type: String
    },
    errno:{
      type: Number,
      required: true
    }
  },
  computed:{
    isEmpty(){
      return this.data.length === 0
    }
  },
  unmounted() {
    this.$emit("flush")
  }
}
</script>

<style scoped>
.worker{
  background-color: wheat;
  text-align: center;
  height: 100%;
}

.worker #worker{
  width: 100%;
  color: darkgreen;
}
</style>