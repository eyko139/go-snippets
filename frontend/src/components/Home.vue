<script lang="ts" setup>
import { onMounted, defineModel, type Ref, ref, watch } from "vue";
import type { Snippet } from "@/types/snippets";
import { getHomeSnippets, getWebsocketClients, connectWs } from "@/api/snippets";
import Button from 'primevue/button';
import Textarea from 'primevue/textarea';

let snippets: Ref<Snippet[]> = ref([]);
let socketConnections: Ref<string[]> = ref([]);
let selectedConnection: Ref<string> = ref("");
let ws: Ref<WebSocket | undefined> = ref();
let wsMessageOut: Ref<string> = ref("");

let message = defineModel<string>("");


onMounted(() => {


  ws.value = connectWs()

  getHomeSnippets().then(res => {
    snippets.value = res;
  })

})

const onClickConnection = (connection: string) => {
  selectedConnection.value = connection;
}

const sendMessage = () => {
  if (!ws.value || !message.value) {
    return;
  }
  const wsMessage = JSON.stringify({
    recipient: selectedConnection.value,
    message: message.value
  })
  ws.value.send(wsMessage);
}

watch(ws, async () => {
  if (!ws.value) {
    return;
  }
  ws.value.onopen = () => {
    getWebsocketClients().then(res => {
      socketConnections.value = res;
    })
  }

  ws.value.onmessage = (event) => {
    wsMessageOut.value = event.data
  }
})

</script>

<template>
  <h1>Snippets</h1>
  <ul v-if="snippets">
    <li v-for="snippet in snippets">
      {{ snippet.title }}
    </li>
  </ul>
  <ul>
    <li :class="{ active: selectedConnection === connection }" v-for="connection in socketConnections"
      @click="onClickConnection(connection)">
      {{ connection }}
    </li>
  </ul>
  <div>
    Message:
    {{ wsMessageOut }}
  </div>
  <Textarea v-model="message" />
  <Button @click="sendMessage" type="submit">Press</Button>
</template>

<style>
.active {
  border: 1px solid green;
}
</style>
