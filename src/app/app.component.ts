import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { SocketService } from "./socket.service";

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css'],
    standalone: true,
    imports: [CommonModule, FormsModule],
    providers: [SocketService]
})
export class AppComponent implements OnInit, OnDestroy {

    public messages: Array<any>;
    public chatBox: string;

    public constructor(private socket: SocketService) {
        this.messages = [];
        this.chatBox = "";
    }

    public ngOnInit() {
      this.socket.getEventListener().subscribe(event => {
          if(event.type == "message") {
              let data = event.data.content;
              if(event.data.sender) {
                  data = event.data.sender + ": " + data;
              }
              this.messages.push(data);
          }
          if(event.type == "close") {
              this.messages.push("/The socket connection has been closed");
          }
          if(event.type == "open") {
              this.messages.push("/The socket connection has been established");
          }
      });
  }

    public ngOnDestroy() {
        this.socket.close();
    }

    public send() {
        if(this.chatBox) {
            this.socket.send(this.chatBox);
            this.chatBox = "";
        }
    }

    public isSystemMessage(message: string) {
        return message.startsWith("/") ? "<strong>" + message.substring(1) + "</strong>" : message;
    }

    

}