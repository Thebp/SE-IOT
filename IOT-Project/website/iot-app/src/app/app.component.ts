import { Component, EventEmitter, Output } from '@angular/core';
import {WebsocketService} from './websocket.service'
import {WsadapterService} from './wsadapter.service'
import {ColorSliderComponent} from './color-picker/color-slider/color-slider.component'

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
  providers: [WebsocketService, WsadapterService]
})
export class AppComponent { 
  title = 'Pycom Information Dashboard';

  constructor(private adapterService: WsadapterService){
    adapterService.messages.subscribe(msg =>{
      console.log("Response from websocket: " + msg);
    });
  }

  //private ctx: CanvasRenderingContext2D;
  
 
  private message = {
    author: "Christian",
    message: "Daniel sover"
  }
/*   getColorAtPosition(x: number, y: number) {
    const imageData = this.ctx.getImageData(x, y, 1, 1).data;
    return 'rgba(' + imageData[0] + ',' + imageData[1] + ',' + imageData[2] + ',1)';
  } */
    
 /*  emitColor(x: number, y: number) {
    const rgbaColor = this.getColorAtPosition(x, y);
    this.color.emit(rgbaColor);
  }
  private mousedown: boolean = false;
  private selectedHeight: number;

  onMouseDown(evt: MouseEvent){
    this.mousedown = true;
    this.selectedHeight = evt.offsetY;
    this.emitColor(evt.offsetX, evt.offsetY);
  }

  onMouseMove(evt: MouseEvent) {
    if (this.mousedown) {
      this.selectedHeight = evt.offsetY;
      this.emitColor(evt.offsetX, evt.offsetY);
    }
  }*/
  //private ctx: CanvasRenderingContext2D = this.canvasRef.nativeElement.getContext("2d");
  sendMsg(){
    console.log("new message from client to Websocket " , this.message);
    this.adapterService.messages.next(this.message);
    this.message.message = "this.ctx.getImageData(x, y,1, 1).data";

  }
}
