import { useEffect, useRef, useState } from 'react';

const WS_URL = import.meta.env.VITE_API_BASE_URL?.replace('http', 'ws') || 'ws://localhost:8088';

interface WebSocketMessage {
  type: string;
  chat_id?: string;
  message_id?: string;
  sender_id?: string;
  content?: string;
  created_at?: string;
  data?: any;
}

type MessageHandler = (message: WebSocketMessage) => void;

export const useWebSocket = () => {
  const [isConnected, setIsConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const handlersRef = useRef<Map<string, MessageHandler[]>>(new Map());
  const reconnectTimeoutRef = useRef<number | undefined>(undefined); 
  const shouldReconnectRef = useRef(true);

  const connect = () => {
    const token = localStorage.getItem('token');
    if (!token) return;

    if (wsRef.current?.readyState === WebSocket.OPEN) {
      console.log('WebSocket already connected');
      setIsConnected(true);
      return;
    }

    if (wsRef.current?.readyState === WebSocket.CONNECTING) {
      console.log('WebSocket is connecting, skipping...');
      return;
    }

    console.log('Connecting to WebSocket...');
    const ws = new WebSocket(`${WS_URL}/ws?token=${token}`);

    ws.onopen = () => {
      console.log('✅ WebSocket connected');
      setIsConnected(true);
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
        reconnectTimeoutRef.current = undefined;
      }
    };

    ws.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data);
        console.log('📨 WebSocket message received:', message);

        const handlers = handlersRef.current.get(message.type) || [];
        handlers.forEach(handler => handler(message));

        const wildcardHandlers = handlersRef.current.get('*') || [];
        wildcardHandlers.forEach(handler => handler(message));
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };

    ws.onerror = (error) => {
      console.error('❌ WebSocket error:', error);
    };

    ws.onclose = () => {
      console.log('🔌 WebSocket disconnected');
      setIsConnected(false);
      wsRef.current = null;

      if (shouldReconnectRef.current) {
        console.log('⏳ Reconnecting in 3 seconds...');
        reconnectTimeoutRef.current = window.setTimeout(() => { 
          if (shouldReconnectRef.current) {
            connect();
          }
        }, 3000);
      }
    };

    wsRef.current = ws;
  };

  const disconnect = () => {
    shouldReconnectRef.current = false;
    
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = undefined;
    }
    if (wsRef.current) {
      wsRef.current.close(1000, 'Component unmounting');
      wsRef.current = null;
    }
    setIsConnected(false);
  };

  const send = (message: WebSocketMessage) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      console.log('📤 Sending:', message); 
      wsRef.current.send(JSON.stringify(message));
    } else {
      console.warn('⚠️ WebSocket is not connected, state:', wsRef.current?.readyState);
    }
  };

  const on = (type: string, handler: MessageHandler) => {
    const handlers = handlersRef.current.get(type) || [];
    handlers.push(handler);
    handlersRef.current.set(type, handlers);

    return () => {
      const currentHandlers = handlersRef.current.get(type) || [];
      const index = currentHandlers.indexOf(handler);
      if (index > -1) {
        currentHandlers.splice(index, 1);
        handlersRef.current.set(type, currentHandlers);
      }
    };
  };

  useEffect(() => {
    shouldReconnectRef.current = true;
    
    const timer = setTimeout(() => {
      connect();
    }, 100);
    
    return () => {
      clearTimeout(timer);
      disconnect();
    };
  }, []); 

  return { isConnected, send, on };
};