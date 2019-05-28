import { Injectable } from '@angular/core';
import {Subject} from 'rxjs';
import {WebsocketService} from './websocket.service'
import {map} from 'rxjs/operators'

const CHAT_URL = "ws://echo.websocket.org/";

export interface Message{
  author: string;
  message: string;
}
@Injectable({
  providedIn: 'root'
})
export class WsadapterService {
public messages: Subject<Message>;

  constructor(wsService : WebsocketService) {
    this.messages = <Subject<Message>>wsService.connect(CHAT_URL).pipe(map(
      (response: MessageEvent): Message => {
        let data = JSON.parse(response.data);
        return {
          author: data.author,
          message: data.message
        }
      }
    ))
   }
}
