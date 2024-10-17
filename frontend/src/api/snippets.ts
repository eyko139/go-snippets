import type {Snippet} from "@/types/snippets";

export const getHomeSnippets = async (): Promise<Snippet[]> => {
  return await fetch(import.meta.env.VITE_BACKEND_URL).then(res => res.json());
}

export const getWebsocketClients = async (): Promise<string[]> => {
  return await fetch(`${import.meta.env.VITE_BACKEND_URL}/wsinfo`).then(res => res.json());
}

export const connectWs = (): WebSocket  => {
  return new WebSocket('ws://localhost:4000/ws')
}
