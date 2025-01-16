import { Injectable, EventEmitter } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class SocketService {

    private socket: WebSocket;
    private listener: EventEmitter<any> = new EventEmitter();

    public constructor() {
        this.socket = new WebSocket("ws://localhost:12345/ws");
        
        this.socket.onopen = (event: Event) => {
            this.listener.emit({"type": "open", "data": event});
        };

        this.socket.onclose = (event: CloseEvent) => {
            this.listener.emit({"type": "close", "data": event});
        };

        this.socket.onmessage = (event: MessageEvent) => {
            try {
                const parsedData = JSON.parse(event.data);
                this.listener.emit({"type": "message", "data": parsedData});
            } catch (e) {
                console.error('Error parsing message:', e);
            }
        };
    }

    public send(data: string) {
        if (this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(data);
        } else {
            console.warn('WebSocket is not open. Message not sent.');
        }
    }

    public close() {
        if (this.socket) {
            this.socket.close();
        }
    }

    public getEventListener() {
        return this.listener;
    }

}